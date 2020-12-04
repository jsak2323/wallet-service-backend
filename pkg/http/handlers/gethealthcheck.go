package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

type GetHealthCheckHandlerResponseMap map[string][]GetHealthCheckRes

type GetHealthCheckService struct {
    moduleServices  *modules.ModuleServiceMap
    healthCheckRepo hc.HealthCheckRepository
}

func NewGetHealthCheckService(moduleServices *modules.ModuleServiceMap, healthCheckRepo hc.HealthCheckRepository) *GetHealthCheckService {
    return &GetHealthCheckService{
        moduleServices,
        healthCheckRepo,
    }
}

func (ghcs *GetHealthCheckService) GetHealthCheckHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    isGetAll := symbol == ""

    RES := make(GetHealthCheckHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetHealthCheckHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetHealthCheckHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    ghcs.InvokeGetHealthCheck(&RES, symbol)
    
    // handle success response
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetHealthCheckHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (ghcs *GetHealthCheckService) InvokeGetHealthCheck(RES *GetHealthCheckHandlerResponseMap, symbol string) {
    // fetch healthcheck data from db
    if symbol != "" { // get by rpc config id
        SYMBOL := strings.ToUpper(symbol)

        for _, rpcConfig := range config.CURR[SYMBOL].RpcConfigs {
            _RES := GetHealthCheckRes{}

            if rpcConfig.Id == 0 { // currency not found
                _RES.Error = "Invalid currency."
                return 
            } 
            _RES.RpcConfig = RpcConfigResDetail{
                RpcConfigId         : rpcConfig.Id,
                Symbol              : SYMBOL,
                Name                : rpcConfig.Name,
                Host                : rpcConfig.Host,
                Type                : rpcConfig.Type,
                NodeVersion         : rpcConfig.NodeVersion,
                NodeLastUpdated     : rpcConfig.NodeLastUpdated,
            }

            healthCheck, err := ghcs.healthCheckRepo.GetByRpcConfigId(rpcConfig.Id)
            if err != nil {
                logger.ErrorLog(" - InvokeGetHealthCheck "+SYMBOL+" GetByRpcConfigId(rpcConfig.Id) err: "+err.Error())
                _RES.Error = err.Error()
                return
            }

            _RES.HealthCheck = healthCheck
            _RES.IsHealthCheckEnabled = rpcConfig.IsHealthCheckEnabled
            
            (*RES)[SYMBOL] = append((*RES)[SYMBOL], _RES)
        }

    } else { // get all

        healthChecks, err := ghcs.healthCheckRepo.GetAllWithRpcConfig()
        if err != nil {
            logger.ErrorLog(" - InvokeGetHealthCheck GetAllWithRpcConfig err: "+err.Error())
            return
        }

        for _, healthCheck := range healthChecks {
            SYMBOL := config.SYMBOLS[healthCheck.RpcConfig.CurrencyId]

            _RES := GetHealthCheckRes{
                RpcConfig: RpcConfigResDetail{
                    RpcConfigId         : healthCheck.RpcConfig.Id,
                    Symbol              : SYMBOL,
                    Name                : healthCheck.RpcConfig.Name,
                    Host                : healthCheck.RpcConfig.Host,
                    Type                : healthCheck.RpcConfig.Type,
                    NodeVersion         : healthCheck.RpcConfig.NodeVersion,
                    NodeLastUpdated     : healthCheck.RpcConfig.NodeLastUpdated,
                },
                HealthCheck: hc.HealthCheck{
                    Id           : healthCheck.Id,
                    RpcConfigId  : healthCheck.RpcConfig.Id,
                    BlockCount   : healthCheck.BlockCount,
                    BlockDiff    : healthCheck.BlockDiff,
                    IsHealthy    : healthCheck.IsHealthy,
                    LastUpdated  : healthCheck.LastUpdated,
                },
                IsHealthCheckEnabled: healthCheck.RpcConfig.IsHealthCheckEnabled
            }

            (*RES)[SYMBOL] = append((*RES)[SYMBOL], _RES)
        }
    }
}


