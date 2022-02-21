package coldwallet

import (
	"encoding/json"
	"net/http"

	handlers "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) UpdateBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq handlers.UpdateBalanceReq
	var RES handlers.StandardRes
	var err error
	var ctx = req.Context()

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			RES.Success = true
			RES.Message = "Cold balance successfully updated"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.ColdWallet
	if err = service.UpdateBalance(ctx, updateReq.Id, updateReq.Balance); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
