package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
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
    tokenType       := vars["token_type"]
    date            := vars["date"]
    rpcConfigType   := vars["rpcconfigtype"]

    SYMBOL := strings.ToUpper(symbol)
    TOKENTYPE := strings.ToUpper(tokenType)
    logger.InfoLog(" - GetLogHandler For symbol: "+SYMBOL+", date: "+date+", type: "+rpcConfigType+", Requesting ...", req) 

    currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, TOKENTYPE)
    if err != nil {
        logger.ErrorLog(" - GetLogHandler config.GetCurrencyBySymbol err: "+err.Error())
        return
    }

    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(currencyConfig.Id, rpcConfigType)
    if err != nil {
        logger.ErrorLog(" - GetLogHandler util.GetRpcConfigByType err: "+err.Error())
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


