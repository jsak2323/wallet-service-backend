package walletrpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetHealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	isGetAll := symbol == ""
	ctx := req.Context()
	status := http.StatusOK

	if isGetAll {
		logger.InfoLog(" - GetHealthCheckHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - GetHealthCheckHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	RES, err := service.InvokeGetHealthCheck(ctx, symbol, tokenType)
	if err != nil {
		status = http.StatusInternalServerError
		logger.ErrorLog(errs.Logged(err), ctx)
	} else {
		logger.InfoLog(" - GetHealthCheckHandler Success.", req)
	}

	// handle success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(RES)
}
