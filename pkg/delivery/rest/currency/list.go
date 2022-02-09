package currency

import (
	"encoding/json"
	"net/http"

	currencyHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES currencyHandler.ListRes
		err error
		ctx = req.Context()
	)

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

	logger.InfoLog(" - currency.ListHandler For all symbols, Requesting ...", req)

	service := re.svc.Currency
	if RES, err = service.ListCurrency(ctx); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllCurrencyConfig)
		return
	}
}
