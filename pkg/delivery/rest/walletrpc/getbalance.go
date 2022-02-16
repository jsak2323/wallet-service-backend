package walletrpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	isGetAll := symbol == ""
	ctx := req.Context()

	if isGetAll {
		logger.InfoLog(" - GetBalanceHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - GetBalanceHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	RES := service.InvokeGetBalance(ctx, symbol, tokenType)

	// handle success response
	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}
