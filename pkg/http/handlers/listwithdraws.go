package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type ListWithdrawsHandlerResponseMap map[string]map[string][]ListWithdrawsRes // map by symbol, token_type

type ListWithdrawsService struct {
	moduleServices *modules.ModuleServiceMap
}

func NewListWithdrawsService(moduleServices *modules.ModuleServiceMap) *ListWithdrawsService {
	return &ListWithdrawsService{
		moduleServices,
	}
}

func (lts *ListWithdrawsService) ListWithdrawsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	limit := vars["limit"]
	isGetAll := symbol == ""
	ctx := req.Context()

	RES := make(ListWithdrawsHandlerResponseMap)

	if isGetAll {
		logger.InfoLog(" - ListWithdrawsHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - ListWithdrawsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	limitInt, _ := strconv.Atoi(limit)
	lts.InvokeListWithdraws(ctx, &RES, symbol, tokenType, limitInt)

	// handle success response
	logger.InfoLog(" - ListWithdrawsHandler Success. Symbol: "+symbol, req)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}

func (lts *ListWithdrawsService) InvokeListWithdraws(ctx context.Context, RES *ListWithdrawsHandlerResponseMap, symbol, tokenType string, limit int) {
	rpcConfigCount := 0
	resChannel := make(chan ListWithdrawsRes)

	for _, curr := range config.CURRRPC {
		SYMBOL := strings.ToUpper(curr.Config.Symbol)
		TOKENTYPE := strings.ToUpper(curr.Config.TokenType)

		// if symbol is defined, only get for that symbol
		if symbol != "" && strings.ToUpper(symbol) != SYMBOL {
			continue
		}
		if tokenType != "" && strings.ToUpper(tokenType) != TOKENTYPE {
			continue
		}

		for _, rpcConfig := range curr.RpcConfigs {
			rpcConfigCount++

			_RES := ListWithdrawsRes{
				RpcConfig: RpcConfigResDetail{
					RpcConfigId:          rpcConfig.Id,
					Symbol:               curr.Config.Symbol,
					TokenType:            curr.Config.TokenType,
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
				module, err := lts.moduleServices.GetModule(currencyConfig.Id)
				if err != nil {
					logger.ErrorLog(" - ListWithdrawsHandler lts.moduleServices.GetModule err: "+err.Error(), ctx)
					_RES.Error = err.Error()
					return
				}

				rpcRes, err := module.ListWithdraws(rpcConfig, limit)
				if err != nil {
					logger.ErrorLog(" - ListWithdrawsHandler (*lts.moduleServices)[SYMBOL].ListWithdraws(rpcConfig, limit) Error: "+err.Error(), ctx)
					_RES.Error = err.Error()

				} else {
					logger.Log(" - InvokeListWithdraws Symbol: "+currencyConfig.Symbol+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host, ctx)
					_RES.Withdraws = rpcRes.Withdraws
					_RES.Error = rpcRes.Error
				}

				resChannel <- _RES

			}(curr.Config, rpcConfig)
		}
	}

	i := 0
	for res := range resChannel {
		i++
		_, ok := (*RES)[res.RpcConfig.Symbol]
		if !ok {
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]ListWithdrawsRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)

		if i >= rpcConfigCount {
			close(resChannel)
		}
	}
}
