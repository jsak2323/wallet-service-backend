package walletrpc

// import (
// 	"io"
// 	"net/http"
// 	"strings"

// 	"github.com/gorilla/mux"

// 	"github.com/btcid/wallet-services-backend-go/cmd/config"
// 	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
// 	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
// 	"github.com/btcid/wallet-services-backend-go/pkg/modules"
// )

// type GetLogService struct {
// 	moduleServices *modules.ModuleServiceMap
// }

// func NewGetLogService(moduleServices *modules.ModuleServiceMap) *GetLogService {
// 	return &GetLogService{
// 		moduleServices,
// 	}
// }

// func (gls *GetLogService) GetLogHandler(w http.ResponseWriter, req *http.Request) {
// 	// define request params
// 	var (
// 		vars          = mux.Vars(req)
// 		symbol        = vars["symbol"]
// 		tokenType     = vars["token_type"]
// 		date          = vars["date"]
// 		rpcConfigType = vars["rpcconfigtype"]

// 		SYMBOL                = strings.ToUpper(symbol)
// 		TOKENTYPE             = strings.ToUpper(tokenType)
// 		errField  *errs.Error = nil
// 		ctx                   = req.Context()
// 	)
// 	logger.InfoLog(" - GetLogHandler For symbol: "+SYMBOL+", date: "+date+", type: "+rpcConfigType+", Requesting ...", req)

// 	defer func() {
// 		if errField != nil {
// 			logger.ErrorLog(errs.Logged(errField), ctx)
// 		}
// 	}()

// 	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, TOKENTYPE)
// 	if err != nil {
// 		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
// 		return
// 	}

// 	// define rpc config
// 	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, rpcConfigType)
// 	if err != nil {
// 		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
// 		return
// 	}

// 	// get log file
// 	res, err := http.Get("http://" + rpcConfig.Host + ":" + rpcConfig.Port + "/log/" + date)
// 	if err != nil {
// 		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetLogFile)
// 		return
// 	}
// 	defer res.Body.Close()

// 	// serve log file
// 	w.Header().Set("Content-Disposition", "attachment; filename=app.log")
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Length", res.Header.Get("Content-Length"))

// 	io.Copy(w, res.Body)
// }
