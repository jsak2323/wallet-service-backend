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

func (re *Rest) MaintenanceListHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	action := vars["action"]
	value := strings.ToUpper(vars["value"])
	ctx := req.Context()

	// define response object
	RES := handlers.StandardRes{}

	// define response handler
	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	service := re.svc.WalletRpc

	err := service.ListMaintenance(ctx, action, value)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
