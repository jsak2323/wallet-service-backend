package rpcconfig

import (
	"encoding/json"
	"net/http"

	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) CreateRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpReq handlerRpcConfig.RpcConfigRpcMethodReq
		RES   handlerRpcConfig.StandardRes
		err   error
		ctx   = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
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
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.RpcConfig
	if err = service.CreateRpcMethod(ctx, rpReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
