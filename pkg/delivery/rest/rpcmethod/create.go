package rpcmethod

import (
	"encoding/json"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	handlerRpcMethod "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcMethod domain.RpcMethod
		RES       handlerRpcMethod.StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcmethod.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Method successfully created"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.RpcMethod
	if err = service.Create(ctx, rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfigRPCMethod)
		return
	}
}
