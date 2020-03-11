package handlers

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/btcid/wallet-services-backend/cmd/config"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    ethservice "github.com/btcid/wallet-services-backend/pkg/modules/eth"

    "github.com/gorilla/mux"
)

type GetBlockCountHandlerResponse map[string][]GetBlockCountRes

func GetBlockCountHandler(w http.ResponseWriter, r *http.Request) { 
    vars := mux.Vars(r)
    symbol := vars["symbol"]

    RES := make(GetBlockCountHandlerResponse)

    var handleSuccess = func() {
        resJson, _ := json.Marshal(RES)
        logger.InfoLog("GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), r)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(RES)
    }

    symbolText := "Symbol: "+symbol
    if symbol == "" { symbolText = "For all symbol" }
    logger.InfoLog("GetBlockCountHandler "+symbolText+", Requesting ...", r) 

    switch symbol { 
        case "eth" :
            getBlockCountEth(&RES)
            handleSuccess()

        default : // get all
            fmt.Println("config.CURR: ")
            ppJson , _ := json.MarshalIndent(config.CURR, "", "\t");
            fmt.Println()
            fmt.Print(string(ppJson))
            fmt.Println()
    }
}

func handleError(err error, funcName string) {
    errMsg := "GetBlockCountHandler "+funcName+" Error: "+err.Error()
    logger.ErrorLog(errMsg)
}

func getBlockCountEth(RES *GetBlockCountHandlerResponse) {
    (*RES)["ETH"] = make([]GetBlockCountRes, 0)
    for _, ethRpcConfig := range config.CURR["ETH"].RpcConfigs {
        res, err := ethservice.GetBlockCount(ethRpcConfig)
        if err != nil { handleError(err, "ethservice.GetBlockCount(ethRpcConfig)") }

        (*RES)["ETH"] = append((*RES)["ETH"], GetBlockCountRes{
            Host    : ethRpcConfig.Host,
            Type    : ethRpcConfig.Type,
            Blocks  : res.Blocks,
        })
    }
}
