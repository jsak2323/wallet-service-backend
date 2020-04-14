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

type GetBlockCountHandlerResponseMap map[string][]GetBlockCountRes

type GetBlockCountService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetBlockCountService(moduleServices *modules.ModuleServiceMap) *GetBlockCountService {
    return &GetBlockCountService{
        moduleServices,
    }
}

func (gbcs *GetBlockCountService) GetBlockCountHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    isGetAll := symbol == ""

    RES := make(GetBlockCountHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBlockCountHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBlockCount(&RES, symbol)
    
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBlockCountService) InvokeGetBlockCount(RES *GetBlockCountHandlerResponseMap, symbol string) {
    var wg sync.WaitGroup
    rpcConfigCount := 0
    resChannel := make(chan GetBlockCountRes)

    for confKey, currConfig := range config.CURR {
        confKey = strings.ToUpper(confKey)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != confKey { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            wg.Add(1)
            rpcConfigCount++
            wg.Done()

            go func(confKey string, rpcConfig rc.RpcConfig) {
                rpcRes, err := (*gbcs.moduleServices)[confKey].GetBlockCount(rpcConfig)
                if err != nil { 
                    logger.ErrorLog(" - GetBlockCountHandler (*gbcs.moduleServices)[confKey].GetBlockCount(rpcConfig) Error: "+err.Error())
                }

                logger.Log(" - InvokeGetBlockCount Symbol: "+confKey+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Blocks: "+rpcRes.Blocks) 
                resChannel <- GetBlockCountRes{
                    RpcConfig: RpcConfigResDetail{
                        RpcConfigId         : rpcConfig.Id,
                        Symbol              : confKey,
                        Name                : rpcConfig.Name,
                        Host                : rpcConfig.Host,
                        Type                : rpcConfig.Type,
                        NodeVersion         : rpcConfig.NodeVersion,
                        NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    },
                    Blocks: rpcRes.Blocks,
                }
            }(confKey, rpcConfig)
        }
    }

    wg.Wait()
    i := 0
    for res := range resChannel {
        i++
        (*RES)[res.Symbol] = append((*RES)[res.Symbol], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}
