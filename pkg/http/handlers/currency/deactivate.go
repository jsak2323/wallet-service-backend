package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *CurrencyConfigService) DeactivateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
		RES    StandardRes
		err    error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- currency.DeactivateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully deactivated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - currency.DeactivateHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.ccRepo.ToggleActive(userId, false); err != nil {
		logger.ErrorLog(" - currency.DeactivateHandler svc.ccRepo.ToggleActive err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
