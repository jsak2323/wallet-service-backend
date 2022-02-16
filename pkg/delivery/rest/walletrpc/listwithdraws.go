package walletrpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) ListWithdrawsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	limitInt, _ := strconv.Atoi(vars["limit"])
	isGetAll := symbol == ""
	ctx := req.Context()
	status := http.StatusOK

	if isGetAll {
		logger.InfoLog(" - ListWithdrawsHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - ListWithdrawsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	res, err := service.InvokeListWithdraws(ctx, symbol, tokenType, limitInt)
	if err != nil {
		status = http.StatusInternalServerError
		logger.ErrorLog(errs.Logged(err), ctx)
	} else {
		logger.InfoLog(" - ListWithdrawsHandler Success. Symbol: "+symbol, req)
	}

	// handle success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
