package rpcresponse

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcResponseService) GetByRpcMethodIdHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error
	var reqRpcMethodId int

	vars := mux.Vars(req)
	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
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
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if RES.RpcResponses, err = s.rrsRepo.GetByRpcMethodId(reqRpcMethodId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedGetRPCResponseByRPCMethodID.Title})
		return
	}
}
