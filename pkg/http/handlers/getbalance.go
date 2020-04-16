package handlers

import (
    "sync"
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

type GetBalanceHandlerResponseMap map[string][]GetBalanceRes

type GetBalanceService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetBalanceService(moduleServices *modules.ModuleServiceMap) *GetBalanceService {
    return &GetBalanceService{
        moduleServices,
    }
}

func (gbcs *GetBalanceService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    isGetAll := symbol == ""

    RES := make(GetBalanceHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBalanceHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBalanceHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBalance(&RES, symbol)

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBalanceService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol string) {
    var wg sync.WaitGroup
    rpcConfigCount := 0
    resChannel := make(chan GetBalanceRes)

    for confKey, currConfig := range config.CURR {
        confKey = strings.ToUpper(confKey)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != confKey { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            wg.Add(1)
            rpcConfigCount++
            wg.Done()

            go func(confKey string, rpcConfig rc.RpcConfig) {
                rpcRes, err := (*gbcs.moduleServices)[confKey].GetBalance(rpcConfig)
                if err != nil { 
                    logger.ErrorLog(" - GetBalanceHandler (*gbcs.moduleServices)[confKey].GetBalance(rpcConfig) Error: "+err.Error())
                }

                logger.Log(" - InvokeGetBalance Symbol: "+confKey+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance) 
                resChannel <- GetBalanceRes{
                    RpcConfig: RpcConfigResDetail{                    
                        RpcConfigId         : rpcConfig.Id,
                        Symbol              : confKey,
                        Name                : rpcConfig.Name,
                        Host                : rpcConfig.Host,
                        Type                : rpcConfig.Type,
                        NodeVersion         : rpcConfig.NodeVersion,
                        NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    },
                    Balance: rpcRes.Balance,
                }
            }(confKey, rpcConfig)
        }
    }

    wg.Wait()
    i := 0
    for res := range resChannel {
        i++
        (*RES)[res.RpcConfig.Symbol] = append((*RES)[res.RpcConfig.Symbol], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}