package currency

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) CreateRpcHandler(w http.ResponseWriter, req *http.Request) {
	var (
		cRreq CurrencyRpcReq
		RES   StandardRes
		err   error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
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

	if err = cRreq.validate(); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.crRepo.Create(cRreq.CurrencyId, cRreq.RpcId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateCurrencyRPC)
		return
	}
}
