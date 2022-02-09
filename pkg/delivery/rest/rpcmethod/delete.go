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

func (re *Rest) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES             handlerRpcMethod.StandardRes
		err             error
		id, RpcConfigId int
		ctx             = req.Context()
	)

	vars := mux.Vars(req)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcmethod.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully deleted"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.DeleteHandler, Requesting ...", req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if RpcConfigId, err = strconv.Atoi(vars["rpc_config_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcMethod
	if err = service.Delete(ctx, id, RpcConfigId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
