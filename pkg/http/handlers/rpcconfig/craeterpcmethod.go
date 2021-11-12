package rpcconfig

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RpcConfigService) CreateRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpReq   RpcConfigRpcMethodReq
		RES   	StandardRes
		err   	error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "RPC Method successfully added to RPC Config"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&rpReq); err != nil {
		logger.ErrorLog(" - rpcconfig.CreateRpcMethodHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if !rpReq.valid() {
		logger.ErrorLog(" - rpcconfig.CreateRpcMethodHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if err = svc.rcrmRepo.Create(rpReq.RpcConfigId, rpReq.RpcMethodId); err != nil {
		logger.ErrorLog(" - rpcconfig.CreateRpcMethodHandler svc.rcrmRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}