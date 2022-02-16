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

type GetBlockCountHandlerResponseMap map[string]map[string][]handlers.GetBlockCountRes

func (s *walletRpcService) InvokeGetBlockCount(ctx context.Context, symbol, tokenType string) (RES *GetBlockCountHandlerResponseMap, err error) {
	var (
		rpcConfigCount = 0
		resChannel     = make(chan handlers.GetBlockCountRes)
	)

	maintenanceList, err := s.GetMaintenanceList(ctx, s.systemConfigRepo)
	if err != nil {
		return nil, errs.AssignErr(errs.AddTrace(err), errs.FailedGetMaintenanceList)
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

			_RES := handlers.GetBlockCountRes{
				RpcConfig: handlers.RpcConfigResDetail{
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
				module, err := s.moduleServices.GetModule(currencyConfig.Id)
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
			(*RES)[res.RpcConfig.Symbol] = make(map[string][]handlers.GetBlockCountRes)
		}

		(*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)
		if i >= rpcConfigCount {
			close(resChannel)
		}
	}

	return RES, nil
}
