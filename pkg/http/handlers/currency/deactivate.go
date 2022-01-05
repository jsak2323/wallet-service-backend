package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- currency.DeactivateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully deactivated"

			config.LoadCurrencyConfigs()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = svc.ccRepo.ToggleActive(userId, false); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedDeactivateCurrencyConfig.Title})
		return
	}
}
