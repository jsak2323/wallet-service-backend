package walletrpc

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type ListTransactionsHandlerResponseMap map[string]map[string][]handlers.ListTransactionsRes // map by symbol, token_type

func (s *walletRpcService) InvokeListTransactions(ctx context.Context, symbol, tokenType string, limit int) (RES *ListTransactionsHandlerResponseMap, err error) {
	var (
		rpcConfigCount = 0
		resChannel     = make(chan handlers.ListTransactionsRes)
	)

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

			_RES := handlers.ListTransactionsRes{
				RpcConfig: handlers.RpcConfigResDetail{
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
				module, err := s.moduleServices.GetModule(currencyConfig.Id)
				if err != nil {
					err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
					return
				}

				rpcRes, err := module.ListTransactions(ctx, rpcConfig, limit)
				if err != nil {
					err = errs.AssignErr(errs.AddTrace(err), errs.FailedListTransactions)

				} else {
					logger.Log(" - InvokeListTransactions Symbol: "+currencyConfig.Symbol+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host, ctx)
					_RES.Transactions = rpcRes.Transactions
					err = errs.AssignErr(errs.AddTrace(errors.New(rpcRes.Error)), errs.FailedListTransactions)
				}

				resChannel <- _RES

				// return err
			}(curr.Config, rpcConfig)
			// if err != nil {
			// 	return RES, err
			// }
		}
	}

	i := 0
	for res := range resChannel {
		i++
		_, ok := (*RES)[res.RpcConfig.Symbol]
		if !ok {
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]handlers.ListTransactionsRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)

		if i >= rpcConfigCount {
			close(resChannel)
		}
	}

	return RES, err
}
