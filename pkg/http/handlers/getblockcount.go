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

    for SYMBOL, currConfig := range config.CURR {
        SYMBOL = strings.ToUpper(SYMBOL)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            wg.Add(1)
            rpcConfigCount++
            wg.Done()

            go func(SYMBOL string, rpcConfig rc.RpcConfig) {
                rpcRes, err := (*gbcs.moduleServices)[SYMBOL].GetBlockCount(rpcConfig)
                if err != nil { 
                    logger.ErrorLog(" - GetBlockCountHandler (*gbcs.moduleServices)[SYMBOL].GetBlockCount(rpcConfig) Error: "+err.Error())
                }

                logger.Log(" - InvokeGetBlockCount Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Blocks: "+rpcRes.Blocks) 
                resChannel <- GetBlockCountRes{
                    RpcConfig: RpcConfigResDetail{
                        RpcConfigId         : rpcConfig.Id,
                        Symbol              : SYMBOL,
                        Name                : rpcConfig.Name,
                        Host                : rpcConfig.Host,
                        Type                : rpcConfig.Type,
                        NodeVersion         : rpcConfig.NodeVersion,
                        NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    },
                    Blocks  : rpcRes.Blocks,
                    Error   : rpcRes.Error,
                }
            }(SYMBOL, rpcConfig)
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


