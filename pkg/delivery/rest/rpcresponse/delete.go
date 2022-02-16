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

func (re *Rest) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES             rpcResponseHandler.StandardRes
		err             error
		id, RpcMethodId int
		ctx             = req.Context()
	)

	vars := mux.Vars(req)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcresponse.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Response successfully deleted"

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.DeleteHandler, Responseing ...", req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if RpcMethodId, err = strconv.Atoi(vars["rpc_method_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcResponse
	if err = service.Delete(ctx, id, RpcMethodId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
