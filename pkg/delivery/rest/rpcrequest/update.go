package rpcrequest

import (
	"encoding/json"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	handlerRpcRequest "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcRequest domain.UpdateRpcRequest
		RES        handlerRpcRequest.StandardRes
		err        error
		ctx        = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcrequest.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully updated"

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcrequest.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcRequest); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.RpcRequest
	if err = service.Update(ctx, rpcRequest); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
