package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	
	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) DeleteRpcHandler(w http.ResponseWriter, req *http.Request) {
	var (
		currencyId, rpcId int
		RES               StandardRes
		err               error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- currency.DeleteRpcHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc successfully removed"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.DeleteRpcHandler, Requesting ...", req)

	vars := mux.Vars(req)
    if currencyId, err = strconv.Atoi(vars["currency_id"]); err != nil {
		logger.ErrorLog(" - DeleteRpcHandler invalid request")
		RES.Error = "Invalid request currency_id"
	}

    if rpcId, err = strconv.Atoi(vars["rpc_id"]); err != nil {
		logger.ErrorLog(" - DeleteRpcHandler invalid request")
		RES.Error = "Invalid request rpc_id"
	}

	if err = s.crRepo.Delete(currencyId, rpcId); err != nil {
		logger.ErrorLog(" - DeleteRpcHandler svc.rpRepo.Delete err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
