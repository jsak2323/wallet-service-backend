package walletrpc

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetLogHandler(w http.ResponseWriter, req *http.Request) {
	// define request params
	var (
		vars          = mux.Vars(req)
		symbol        = vars["symbol"]
		tokenType     = vars["token_type"]
		date          = vars["date"]
		rpcConfigType = vars["rpcconfigtype"]
		ctx           = req.Context()
	)

	service := re.svc.WalletRpc
	resp, err := service.GetLog(ctx, symbol, tokenType, rpcConfigType, date)
	if err != nil {
		logger.ErrorLog(errs.Logged(errs.AddTrace(err)), ctx)
	}

	// serve log file
	w.Header().Set("Content-Disposition", "attachment; filename=app.log")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

	io.Copy(w, resp.Body)
}
