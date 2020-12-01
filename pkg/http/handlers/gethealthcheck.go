package handlers

import (
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

type GetHealthCheckHandlerResponseMap map[string][]GetHealthCheckRes

type GetHealthCheckService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetHealthCheckService(moduleServices *modules.ModuleServiceMap) *GetHealthCheckService {
    return &GetHealthCheckService{
        moduleServices,
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

    gbcs.InvokeGetHealthCheck(&RES, symbol)
    
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
        _RES   := GetHealthCheckRes{}

        rpcConfig := config.CURR[SYMBOL]
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

        _RES.HealthChecks = append(healthChecks, healthCheck)
        (*RES)[SYMBOL] = append((*RES)[SYMBOL], _RES)


    } else { // get all

        healthChecks, err := ghcs.healthCheckRepo.GetAll()
        if err != nil {
            logger.ErrorLog(" - InvokeGetHealthCheck GetAll err: "+err.Error())
            return
        }

        for SYMBOL, currConfig := range config.CURR {
            SYMBOL = strings.ToUpper(SYMBOL)

            // if symbol is defined, only get for that symbol
            // if symbol != "" && strings.ToUpper(symbol) != SYMBOL { continue }

            for _, rpcConfig := range currConfig.RpcConfigs {
                // rpcConfigCount++

                _RES := GetHealthCheckRes{
                    RpcConfig: RpcConfigResDetail{
                        RpcConfigId         : rpcConfig.Id,
                        Symbol              : SYMBOL,
                        Name                : rpcConfig.Name,
                        Host                : rpcConfig.Host,
                        Type                : rpcConfig.Type,
                        NodeVersion         : rpcConfig.NodeVersion,
                        NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    },
                }



            }
        }

    }
}


