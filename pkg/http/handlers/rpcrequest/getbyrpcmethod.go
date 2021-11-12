package rpcrequest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcRequestService) GetByRpcMethodIdHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error
	var reqRpcMethodId int

	vars := mux.Vars(req)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcrequest.GetByRpcMethodIdHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcrequest.GetByRpcMethodIdHandler, Requesting ...", req)

	if reqRpcMethodId, err = strconv.Atoi(vars["rpc_method_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcMethodIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.RpcRequests, err = s.rrqRepo.GetByRpcMethodId(reqRpcMethodId); err != nil {
		logger.ErrorLog(" -- rpcrequest.GetByRpcMethodIdHandler rrqRepo.GetByRpcMethodId Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
