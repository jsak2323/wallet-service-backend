package cold

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *ColdWalletService) ActivateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		id  int
		RES StandardRes
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

	if err = svc.cbRepo.ToggleActive(ctx, id, true); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedActivatedColdBalance)
		return
	}
}
