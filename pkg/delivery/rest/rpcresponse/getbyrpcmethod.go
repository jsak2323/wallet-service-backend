package rpcresponse

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	rpcResponseHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetByRpcMethodIdHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES            rpcResponseHandler.ListRes
		err            error
		reqRpcMethodId int
		ctx            = req.Context()
	)

	vars := mux.Vars(req)
	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" - rpcresponse.GetByRpcMethodIdHandler, success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcresponse.GetByRpcMethodIdHandler, Requesting ...", req)

	if reqRpcMethodId, err = strconv.Atoi(vars["rpc_method_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcResponse
	if RES.RpcResponses, err = service.GetByRpcMethod(ctx, reqRpcMethodId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
