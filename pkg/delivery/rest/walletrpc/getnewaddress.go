package walletrpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetNewAddressHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	addressType := vars["type"]

	// define response handler

	service := re.svc.WalletRpc
	res, err := service.GetNewAddress(ctx, symbol, tokenType, addressType)
	if err != nil {
		res.Error = errs.AddTrace(err)
		return
	}

	handleResponse := func() {
		resStatus := http.StatusOK
		if res.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(res.Error), ctx)
		} else {
			resJson, _ := json.Marshal(res)
			logger.InfoLog(" - GetNewAddressHandler Success. Symbol: "+strings.ToUpper(symbol)+", Res: "+string(resJson), req)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(res)
	}
	defer handleResponse()

}
