package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type GetBlockCountHandlerResponseMap map[string]map[string][]GetBlockCountRes

type GetBlockCountService struct {
	moduleServices   *modules.ModuleServiceMap
	systemConfigRepo sc.Repository
}

func NewGetBlockCountService(
	moduleServices *modules.ModuleServiceMap,
	systemConfigRepo sc.Repository,
) *GetBlockCountService {
	return &GetBlockCountService{
		moduleServices,
		systemConfigRepo,
	}
}

func (gbcs *GetBlockCountService) GetBlockCountHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	isGetAll := symbol == ""
	ctx := req.Context()

	RES := make(GetBlockCountHandlerResponseMap)

	if isGetAll {
		logger.InfoLog(" - GetBlockCountHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	gbcs.InvokeGetBlockCount(ctx, &RES, symbol, tokenType)

	// handle success response
	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBlockCountService) InvokeGetBlockCount(ctx context.Context, RES *GetBlockCountHandlerResponseMap, symbol, tokenType string) {
	var (
		rpcConfigCount             = 0
		resChannel                 = make(chan GetBlockCountRes)
		errField       *errs.Error = nil
		err            error
	)

	defer func() {
		if err != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	maintenanceList, err := GetMaintenanceList(ctx, gbcs.systemConfigRepo)
	if err != nil {
		logger.ErrorLog(errs.Logged(errs.FailedGetMaintenanceList), ctx)
	}

	for _, currRpc := range config.CURRRPC {
		SYMBOL := strings.ToUpper(currRpc.Config.Symbol)
		TOKENTYPE := strings.ToUpper(currRpc.Config.TokenType)

		// if symbol is defined, only get for that symbol
		if symbol != "" && strings.ToUpper(symbol) != SYMBOL {
			continue
		}
		if tokenType != "" && strings.ToUpper(tokenType) != TOKENTYPE {
			continue
		}

		// if not parent coin, skip
		if currRpc.Config.TokenType != cc.MainTokenType {
			continue
		}

		// if maintenance, skip
		if maintenanceList[SYMBOL] {
			continue
		}

		for _, rpcConfig := range currRpc.RpcConfigs {
			rpcConfigCount++

			_RES := GetBlockCountRes{
				RpcConfig: RpcConfigResDetail{
					RpcConfigId:          rpcConfig.Id,
					Symbol:               currRpc.Config.Symbol,
					TokenType:            currRpc.Config.TokenType,
					Name:                 rpcConfig.Name,
					Host:                 rpcConfig.Host,
					Type:                 rpcConfig.Type,
					NodeVersion:          rpcConfig.NodeVersion,
					NodeLastUpdated:      rpcConfig.NodeLastUpdated,
					IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
				},
			}

			// execute concurrent rpc calls
			go func(currencyConfig cc.CurrencyConfig, rpcConfig rc.RpcConfig) {
				module, err := gbcs.moduleServices.GetModule(currencyConfig.Id)
				if err != nil {
					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
					return
				}

				rpcRes, err := module.GetBlockCount(ctx, rpcConfig)
				if err != nil {
					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetBlockCount)
				} else {
					logger.Log(" - InvokeGetBlockCount Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Blocks: "+rpcRes.Blocks, ctx)
					_RES.Blocks = rpcRes.Blocks
					_RES.Error = errs.AssignErr(errs.AddTrace(errors.New(rpcRes.Error)), errs.FailedGetBlockCount)
				}

				resChannel <- _RES

			}(currRpc.Config, rpcConfig)
		}
	}

	i := 0
	for res := range resChannel {
		i++
		_, ok := (*RES)[res.RpcConfig.Symbol]
		if !ok {
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]GetBlockCountRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)
		if i >= rpcConfigCount {
			close(resChannel)
		}
	}
}
