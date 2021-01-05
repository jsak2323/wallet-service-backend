package handlers

import (
    "fmt"
    "strings"
    "net/http"
    // "encoding/json"

    "github.com/gorilla/mux"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type GetLogService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetLogService(moduleServices *modules.ModuleServiceMap) *GetLogService {
    return &GetLogService{
        moduleServices,
    }
}

func (gls *GetLogService) GetLogHandler(w http.ResponseWriter, req *http.Request) { 
    // define response handler
    handleResponse := func() {
        resStatus := http.StatusOK
        // if RES.Error != "" {
        //     resStatus = http.StatusInternalServerError
        // }
        w.WriteHeader(resStatus)
        // json.NewEncoder(w).Encode(RES)
    }
    defer handleResponse()

    // define request params
    vars := mux.Vars(req)
    symbol          := vars["symbol"]
    date            := vars["date"]
    rpcConfigType   := vars["rpcconfigtype"]

    SYMBOL := strings.ToUpper(symbol)
    logger.InfoLog(" - GetLogHandler For symbol: "+SYMBOL+", date: "+date+", type: "+rpcConfigType+", Requesting ...", req) 

    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, rpcConfigType)
    if err != nil {
        logger.ErrorLog(" - GetLogHandler util.GetRpcConfigByType(SYMBOL, rpcConfigType) err: "+err.Error())
        // RES.Error = err.Error()
        return
    }

    fmt.Printf("rpcConfig: %+v\n", rpcConfig)

    // RES.RpcConfig = RpcConfigResDetail{
    //     RpcConfigId             : rpcConfig.Id,
    //     Symbol                  : SYMBOL,
    //     Name                    : rpcConfig.Name,
    //     Host                    : rpcConfig.Host,
    //     Type                    : rpcConfig.Type,
    //     NodeVersion             : rpcConfig.NodeVersion,
    //     NodeLastUpdated         : rpcConfig.NodeLastUpdated,
    //     IsHealthCheckEnabled    : rpcConfig.IsHealthCheckEnabled,
    // }

    // execute rpc call
    filepath, err := (*gls.moduleServices)[SYMBOL].GetLog(rpcConfig, date)
    // if err != nil { 
    //     logger.ErrorLog(" - GetLogHandler (*gls.moduleServices)[SYMBOL].GetLog(rpcConfig, date), Error: "+err.Error())
    //     // RES.Error = err.Error()
    //     return
    // }

    // handle success response
    // RES.AddressType = rpcRes.AddressType
    // RES.Error       = rpcRes.Error
    // resJson, _ := json.Marshal(RES)
    // logger.InfoLog(" - AddressTypeHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}


