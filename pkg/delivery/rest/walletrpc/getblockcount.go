package walletrpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetBlockCountHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	isGetAll := symbol == ""
	ctx := req.Context()
	status := http.StatusOK

	if isGetAll {
		logger.InfoLog(" - GetBlockCountHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
	}

	service := re.svc.WalletRpc
	RES, err := service.InvokeGetBlockCount(ctx, symbol, tokenType)
	if err != nil {
		status = http.StatusInternalServerError
		logger.ErrorLog(errs.Logged(err), ctx)
	} else {
		resJson, _ := json.Marshal(RES)
		logger.InfoLog(" - GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	}

	// handle success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(RES)
}
