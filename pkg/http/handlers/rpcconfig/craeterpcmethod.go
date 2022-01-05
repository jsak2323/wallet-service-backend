package rpcconfig

import (
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RpcConfigService) CreateRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpReq RpcConfigRpcMethodReq
		RES   StandardRes
		err   error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			RES.Success = true
			RES.Message = "RPC Method successfully added to RPC Config"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&rpReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if !rpReq.valid() {
		err = errors.New("Invalid request")
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = svc.rcrmRepo.Create(rpReq.RpcConfigId, rpReq.RpcMethodId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedCreateRPCConfigRPCMethod.Title})
		return
	}
}
