package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) DeleteRpcHandler(w http.ResponseWriter, req *http.Request) {
	var (
		currencyId, rpcId int
		RES               StandardRes
		err               error
		ctx               = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- currency.DeleteRpcHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc successfully removed"

			config.LoadCurrencyConfigs(ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.DeleteRpcHandler, Requesting ...", req)

	vars := mux.Vars(req)
	if currencyId, err = strconv.Atoi(vars["currency_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return

	}

	if rpcId, err = strconv.Atoi(vars["rpc_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.crRepo.Delete(ctx, currencyId, rpcId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteCurrencyRPC)
		return
	}
}
