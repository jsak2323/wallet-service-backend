package walletrpc

import (
	"context"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type GetHealthCheckHandlerResponseMap map[string]map[string][]handlers.GetHealthCheckRes // map by symbol, token_type

func (s *walletRpcService) InvokeGetHealthCheck(ctx context.Context, symbol, tokenType string) (RES *GetHealthCheckHandlerResponseMap, err error) {

	// get maintenance list
	maintenanceList, err := s.GetMaintenanceList(ctx, s.systemConfigRepo)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetMaintenanceList)
	}

	// fetch healthcheck data from db
	if symbol != "" { // get by rpc config id
		SYMBOL := strings.ToUpper(symbol)
		TOKENTYPE := strings.ToUpper(tokenType)

		currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, TOKENTYPE)
		if err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
			return nil, err
		}

		for _, rpcConfig := range config.CURRRPC[currencyConfig.Id].RpcConfigs {
			_RES := handlers.GetHealthCheckRes{}

			if rpcConfig.Id == 0 { // currency not found
				err = errs.AddTrace(errs.InvalidCurrency)
				return nil, err
			}
			_RES.RpcConfig = handlers.RpcConfigResDetail{
				RpcConfigId:          rpcConfig.Id,
				Symbol:               SYMBOL,
				TokenType:            TOKENTYPE,
				Name:                 rpcConfig.Name,
				Host:                 rpcConfig.Host,
				Type:                 rpcConfig.Type,
				NodeVersion:          rpcConfig.NodeVersion,
				NodeLastUpdated:      rpcConfig.NodeLastUpdated,
				IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
			}

			healthCheck, err := s.healthCheckRepo.GetByRpcConfigId(rpcConfig.Id)
			if err != nil {
				err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetByRpcConfigId)
				return nil, err
			}

			_RES.HealthCheck = healthCheck
			_RES.IsMaintenance = maintenanceList[SYMBOL]

			_, ok := (*RES)[SYMBOL]
			if !ok {
				(*RES)[SYMBOL] = make(map[string][]handlers.GetHealthCheckRes)
			}
			(*RES)[SYMBOL][TOKENTYPE] = append((*RES)[SYMBOL][TOKENTYPE], _RES)
		}

	} else { // get all

		for _, curr := range config.CURRRPC {
			for _, rpcConfig := range curr.RpcConfigs {

				healthCheck, err := s.healthCheckRepo.GetByRpcConfigId(rpcConfig.Id)
				if err != nil {
					err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetByRpcConfigId)
					return nil, err
				}

				_RES := handlers.GetHealthCheckRes{
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
					HealthCheck: hc.HealthCheck{
						Id:          healthCheck.Id,
						RpcConfigId: rpcConfig.Id,
						BlockCount:  healthCheck.BlockCount,
						BlockDiff:   healthCheck.BlockDiff,
						IsHealthy:   healthCheck.IsHealthy,
						LastUpdated: healthCheck.LastUpdated,
					},
					IsMaintenance: maintenanceList[curr.Config.Symbol],
				}

				_, ok := (*RES)[curr.Config.Symbol]
				if !ok {
					(*RES)[curr.Config.Symbol] = make(map[string][]handlers.GetHealthCheckRes)
				}
				(*RES)[curr.Config.Symbol][curr.Config.TokenType] = append((*RES)[curr.Config.Symbol][curr.Config.TokenType], _RES)
			}
		}
	}

	return RES, err
}
