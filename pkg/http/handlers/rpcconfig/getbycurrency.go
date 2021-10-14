package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) GetByCurrencyHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error
	var reqCurrencyId int

	vars := mux.Vars(req)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcconfig.GetByCurrencyHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcconfig.GetByCurrencyHandler, Requesting ...", req)

	if reqCurrencyId, err = strconv.Atoi(vars["currency_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByIdHandler strconv.Atoi(" + vars["currency_id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.RpcConfigs, err = s.rcRepo.GetByCurrencyId(reqCurrencyId); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByCurrencyHandler rcRepo.GetByCurrencySymbol Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

}
