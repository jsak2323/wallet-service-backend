package coldwallet

import (
	"encoding/json"
	"net/http"
	"strconv"

	handlers "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (re *Rest) ActivateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		id  int
		RES handlers.StandardRes
		err error
		ctx = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- cold.ActivateHandler Success!", req)

			RES.Success = true
			RES.Message = "Cold Wallet successfully activated"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.ColdWallet
	if err = service.ActivateColdWallet(ctx, id); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
