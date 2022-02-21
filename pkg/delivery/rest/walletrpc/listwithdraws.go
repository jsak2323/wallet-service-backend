package walletrpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/service/walletrpc"
)

func (re *Rest) ListWithdrawsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	limitInt, _ := strconv.Atoi(vars["limit"])
	isGetAll := symbol == ""
	ctx := req.Context()
	status := http.StatusOK
	wdRES := make(walletrpc.ListWithdrawsHandlerResponseMap)

	if isGetAll {
		logger.InfoLog(" - ListWithdrawsHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - ListWithdrawsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	service.InvokeListWithdraws(ctx, &wdRES, symbol, tokenType, limitInt)

	// handle success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(wdRES)
}
