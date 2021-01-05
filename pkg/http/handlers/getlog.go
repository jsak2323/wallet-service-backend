package handlers

import (
    "io"
    "strings"
    "net/http"

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
        return
    }

    // get log file
    res, err := http.Get("http://"+rpcConfig.Host+":"+rpcConfig.Port+"/log/"+date)
    if err != nil { 
        logger.ErrorLog(" - GetLogHandler http.get err: "+err.Error())
        return
    }
    defer res.Body.Close()

    // serve log file
    w.Header().Set("Content-Disposition", "attachment; filename=app.log")
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", res.Header.Get("Content-Length"))

    io.Copy(w, res.Body)
}


