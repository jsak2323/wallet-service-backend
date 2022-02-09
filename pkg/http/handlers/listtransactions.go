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
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type ListTransactionsHandlerResponseMap map[string]map[string][]ListTransactionsRes // map by symbol, token_type

type ListTransactionsService struct {
	moduleServices *modules.ModuleServiceMap
}

func NewListTransactionsService(moduleServices *modules.ModuleServiceMap) *ListTransactionsService {
	return &ListTransactionsService{
		moduleServices,
	}
}

func (lts *ListTransactionsService) ListTransactionsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	limit := vars["limit"]
	isGetAll := symbol == ""
	ctx := req.Context()

	RES := make(ListTransactionsHandlerResponseMap)

	if isGetAll {
		logger.InfoLog(" - ListTransactionsHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - ListTransactionsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	limitInt, _ := strconv.Atoi(limit)
	lts.InvokeListTransactions(ctx, &RES, symbol, tokenType, limitInt)

	// handle success response
	logger.InfoLog(" - ListTransactionsHandler Success. Symbol: "+symbol, req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}

func (lts *ListTransactionsService) InvokeListTransactions(ctx context.Context, RES *ListTransactionsHandlerResponseMap, symbol, tokenType string, limit int) {
	var (
		rpcConfigCount             = 0
		resChannel                 = make(chan ListTransactionsRes)
		errField       *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

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

			_RES := ListTransactionsRes{
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
					errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
					return
				}

				rpcRes, err := module.ListTransactions(ctx, rpcConfig, limit)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.FailedListTransactions)

				} else {
					logger.Log(" - InvokeListTransactions Symbol: "+currencyConfig.Symbol+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host, ctx)
					_RES.Transactions = rpcRes.Transactions
					errField = errs.AssignErr(errs.AddTrace(errors.New(rpcRes.Error)), errs.FailedListTransactions)
					_RES.Error = errField
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
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]ListTransactionsRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)

		if i >= rpcConfigCount {
			close(resChannel)
		}
	}
}
