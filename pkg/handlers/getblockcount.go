package handlers

import (
    "sync"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/cmd/config"
    "github.com/btcid/wallet-services-backend/pkg/modules"
)

type GetBlockCountHandlerResponseMap map[string][]GetBlockCountRes

func GetBlockCountHandler(w http.ResponseWriter, r *http.Request) { 
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
        invokeGetAllBlockCount(&RES)

    } else {
        logger.InfoLog("GetBlockCountHandler For symbol: "+symbol+", Requesting ...", r) 
        invokeGetBlockCount(symbol, &RES)
    }

    handleSuccess()
}

func handleError(err error, funcName string) {
    errMsg := "GetBlockCountHandler "+funcName+" Error: "+err.Error()
    logger.ErrorLog(errMsg)
}

func invokeGetBlockCount(symbol string, RES *GetBlockCountHandlerResponseMap) {
    confKey := strings.ToUpper(symbol)
    ModuleServices := modules.NewModuleServices()

    (*RES)[confKey] = make([]GetBlockCountRes, 0)
    for _, rpcConfig := range config.CURR[confKey].RpcConfigs {
        rpcRes, err := (*ModuleServices)[confKey].GetBlockCount(rpcConfig)
        if err != nil { handleError(err, "invokeGetBlockCount "+confKey+".GetBlockCount(rpcConfig)") }

        (*RES)[confKey] = append((*RES)[confKey], GetBlockCountRes{
            Symbol  : confKey,
            Host    : rpcConfig.Host,
            Type    : rpcConfig.Type,
            Blocks  : rpcRes.Blocks,
        })
    }
}

func invokeGetAllBlockCount(RES *GetBlockCountHandlerResponseMap) {
    var wg sync.WaitGroup
    rpcConfigCount := 0
    resChannel := make(chan GetBlockCountRes)
    ModuleServices := modules.NewModuleServices()

    for confKey, currConfig := range config.CURR {
        confKey = strings.ToUpper(confKey)

        for _, rpcConfig := range currConfig.RpcConfigs {
            wg.Add(1)
            rpcConfigCount++
            wg.Done()

            go func(confKey string, rpcConfig rc.RpcConfig) {
                rpcRes, err := (*ModuleServices)[confKey].GetBlockCount(rpcConfig)
                if err != nil { handleError(err, "invokeGetAllBlockCount "+confKey+".GetBlockCount(rpcConfig)") }

                resChannel <- GetBlockCountRes{
                    Symbol  : confKey,
                    Host    : rpcConfig.Host,
                    Type    : rpcConfig.Type,
                    Blocks  : rpcRes.Blocks,
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
