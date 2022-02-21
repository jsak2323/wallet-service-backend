package userwallet

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	service "github.com/btcid/wallet-services-backend-go/pkg/service/userwallet"
)

func (re *Rest) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := strings.ToUpper(vars["symbol"])
	isGetAll := symbol != ""
	ctx := req.Context()

	RES := make(service.GetBalanceHandlerResponseMap)

	if isGetAll {
		logger.InfoLog(" - userwallet.GetBalanceHandler For all symbols, Requesting ...", req)
	} else {
		logger.InfoLog(" - userwallet.GetBalanceHandler For symbol: "+symbol+", Requesting ...", req)
	}

	service := re.svc.UserWallet
	service.InvokeGetBalance(ctx, &RES, symbol)

	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - userwallet.GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}
