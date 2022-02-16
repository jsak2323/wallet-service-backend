package walletrpc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) AddressTypeHandler(w http.ResponseWriter, req *http.Request) {
	// define response object
	RES := handlers.AddressTypeRes{}
	ctx := req.Context()

	// define response handler
	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	// define request params
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	address := vars["address"]

	service := re.svc.WalletRpc
	RES, err := service.GetAddressType(ctx, symbol, tokenType, address)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - AddressTypeHandler Success. Symbol: "+strings.ToUpper(symbol)+", Res: "+string(resJson), req)
}
