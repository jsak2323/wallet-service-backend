package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *CurrencyConfigService) ActivateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- currency.DeactivateHandler Success!", req)
			
			RES.Success = true
			RES.Message = "Currency successfully activated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - currency.ActivateHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.ccRepo.ToggleActive(id, true); err != nil {
		logger.ErrorLog(" - currency.ActivateHandler svc.ccRepo.ToggleActive err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
