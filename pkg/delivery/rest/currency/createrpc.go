package currency

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	currencyHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) CreateRpcHandler(w http.ResponseWriter, req *http.Request) {
	var (
		cRreq currencyHandler.CurrencyRpcReq
		RES   StandardRes
		err   error
		ctx   = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs(ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.CreateRpcHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&cRreq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.Currency
	if err = service.CreateCurrencyRpc(ctx, cRreq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateCurrencyRPC)
		return
	}
}
