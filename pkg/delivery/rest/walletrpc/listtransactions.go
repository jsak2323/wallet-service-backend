package walletrpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) ListTransactionsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	limit := vars["limit"]
	isGetAll := symbol == ""
	ctx := req.Context()
	limitInt, _ := strconv.Atoi(limit)
	status := http.StatusOK

	if isGetAll {
		logger.InfoLog(" - ListTransactionsHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - ListTransactionsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	res, err := service.InvokeListTransactions(ctx, symbol, tokenType, limitInt)
	if err != nil {
		status = http.StatusInternalServerError
	} else {
		logger.InfoLog(" - ListTransactionsHandler Success. Symbol: "+symbol, req)
	}

	// handle success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
