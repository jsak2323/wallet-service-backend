package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	handlerRpcMethod "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetByRpcConfigIdHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES            handlerRpcMethod.ListRes
		err            error
		reqRpcConfigId int
		ctx            = req.Context()
	)

	vars := mux.Vars(req)
	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" - rpcmethod.ListHandler, success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcmethod.ListHandler, Requesting ...", req)

	if reqRpcConfigId, err = strconv.Atoi(vars["rpc_config_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcMethod
	if RES.RpcMethods, err = service.GetByRpcConfigId(ctx, reqRpcConfigId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
