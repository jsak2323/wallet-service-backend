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

func (gbcs *GetBlockCountService) GetBlockCountHandler(w http.ResponseWriter, r *http.Request) { 
    vars := mux.Vars(r)
    symbol := vars["symbol"]
    isGetAll := symbol == ""

    RES := make(GetBlockCountHandlerResponseMap)

    var handleSuccess = func() {
        resJson, _ := json.Marshal(RES)
        logger.InfoLog("GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), r)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(RES)
    }

    if isGetAll {
        logger.InfoLog("GetBlockCountHandler For all symbols, Requesting ...", r) 
    } else {
        logger.InfoLog("GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", r) 
    }

    gbcs.InvokeGetBlockCount(&RES, symbol)
    handleSuccess()
}

func (gbcs *GetBlockCountService) handleError(err error, funcName string) {
    errMsg := "GetBlockCountHandler "+funcName+" Error: "+err.Error()
    logger.ErrorLog(errMsg)
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
                if err != nil { gbcs.handleError(err, "InvokeGetBlockCount "+confKey+".GetBlockCount(rpcConfig)") }

                logger.Log("InvokeGetBlockCount Symbol: "+confKey+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Blocks: "+rpcRes.Blocks) 
                resChannel <- GetBlockCountRes{
                    RpcConfigId         : rpcConfig.Id,
                    Symbol              : confKey,
                    Host                : rpcConfig.Host,
                    Type                : rpcConfig.Type,
                    NodeVersion         : rpcConfig.NodeVersion,
                    NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    Blocks              : rpcRes.Blocks,
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
