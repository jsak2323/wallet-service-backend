package cold

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *ColdWalletService) DeactivateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		id  int
		RES StandardRes
		err error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- cold.DeactivateHandler Success!", req)

			RES.Success = true
			RES.Message = "Cold Wallet successfully deactivated"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" -- cold.DeactivateHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.cbRepo.ToggleActive(id, false); err != nil {
		logger.ErrorLog(" -- cold.DeactivateHandler svc.cbRepo.ToggleActive err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
