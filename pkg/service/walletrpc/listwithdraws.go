package walletrpc

import (
	"context"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type ListWithdrawsHandlerResponseMap map[string]map[string][]handlers.ListWithdrawsRes // map by symbol, token_type

func (s *walletRpcService) InvokeListWithdraws(ctx context.Context, symbol, tokenType string, limit int) (RES *ListWithdrawsHandlerResponseMap, err error) {
	rpcConfigCount := 0
	resChannel := make(chan handlers.ListWithdrawsRes)

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

			_RES := handlers.ListWithdrawsRes{
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
					err = errs.AddTrace(errs.AddTrace(err))
					return
				}

				rpcRes, err := module.ListWithdraws(ctx, rpcConfig, limit)
				if err != nil {
					err = errs.AddTrace(errs.AddTrace(err))

				} else {
					_RES.Withdraws = rpcRes.Withdraws
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
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]handlers.ListWithdrawsRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)

		if i >= rpcConfigCount {
			close(resChannel)
		}
	}

	return RES, err
}
