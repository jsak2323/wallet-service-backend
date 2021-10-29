package currency

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) CreateRpcHandler(w http.ResponseWriter, req *http.Request) {
	var (
		cRreq CurrencyRpcReq
		RES StandardRes
		err error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- currency.CreateRpcHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.CreateRpcHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&cRreq); err != nil {
		logger.ErrorLog(" -- currency.CreateRpcHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = cRreq.validate(); err != nil {
		logger.ErrorLog(" -- currency.CreateRpcHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.crRepo.Create(cRreq.CurrencyId, cRreq.RpcId); err != nil {
		logger.ErrorLog(" - currency.CreateRpcHandler s.crRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
