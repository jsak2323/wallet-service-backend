package cron

import(
    "fmt"
    "strconv"
    "net/http"
    // "encoding/json"

    h "github.com/btcid/wallet-services-backend/pkg/http/handlers"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // node blockcounts 
    logger.InfoLog("HealthCheckHandler Getting node blockcounts ..." , r)
    gbcRES := make(h.GetBlockCountHandlerResponseMap)
    h.InvokeGetBlockCount(&gbcRES, "", false)
    logger.InfoLog("HealthCheckHandler Getting node blockcounts done. Fetched "+strconv.Itoa(len(gbcRES))+" results." , r)

    // comparation blockcounts
    logger.InfoLog("HealthCheckHandler Getting comparation blockcounts ..." , r)
    cbcRES := make(h.GetBlockCountHandlerResponseMap)
    h.InvokeGetBlockCount(&cbcRES, "", true)
    logger.InfoLog("HealthCheckHandler Getting comparation blockcounts done. Fetched "+strconv.Itoa(len(cbcRES))+" results." , r)

    for resSymbol, resRpcConfigs := range gbcRES {
        for _, resRpcConfig := range resRpcConfigs {        
            nodeBlocks          := resRpcConfig.Blocks
            comparationBlocks   := findComparationResultByHost(resSymbol, resRpcConfig.Host, &cbcRES).Blocks

            fmt.Println("nodeBlocks: "+nodeBlocks)
            fmt.Println("comparationBlocks: "+comparationBlocks)
        }
    }
}

func findComparationResultByHost(symbol string, host string, RES *h.GetBlockCountHandlerResponseMap) h.GetBlockCountRes {
    emptyRes := h.GetBlockCountRes{}
    for resSymbol, resRpcConfigs := range (*RES) {
        for _, resRpcConfig := range resRpcConfigs {
            if symbol == resSymbol && resRpcConfig.Host == host { return resRpcConfig }
        }
    }
    return emptyRes
}
