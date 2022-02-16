package walletrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceHandlerResponseMap map[string][]handlers.GetBalanceRes

func (s *walletRpcService) InvokeGetBalance(ctx context.Context, symbol, tokenType string) (RES *GetBalanceHandlerResponseMap) {
	var (
		rpcConfigCount = 0
		resChannel     = make(chan handlers.GetBalanceRes)
	)

	for _, curr := range config.CURRRPC {
		SYMBOL := strings.ToUpper(curr.Config.Symbol)
		TOKENTYPE := strings.ToUpper(curr.Config.Symbol)

		// if symbol is defined, only get for that symbol
		if symbol != "" && strings.ToUpper(symbol) != SYMBOL && strings.ToUpper(tokenType) != TOKENTYPE {
			continue
		}

		for _, rpcConfig := range curr.RpcConfigs {
			rpcConfigCount++
			_RES := handlers.GetBalanceRes{
				RpcConfig: handlers.RpcConfigResDetail{
					RpcConfigId:          rpcConfig.Id,
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
					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
					return
				}

				rpcRes, err := module.GetBalance(ctx, rpcConfig)
				if err != nil {
					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetBalance)

				} else {
					logger.Log(" -- InvokeGetBalance Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance, ctx)
					_RES.Balance = rpcRes.Balance
					// err = errs.AddTrace(errors.New(rpcRes.Error))
				}

				resChannel <- _RES

			}(curr.Config, rpcConfig)
		}
	}

	i := 0
	for res := range resChannel {
		i++
		(*RES)[res.RpcConfig.Symbol] = append((*RES)[res.RpcConfig.Symbol], res)
		if i >= rpcConfigCount {
			close(resChannel)
		}
	}

	return RES
}
