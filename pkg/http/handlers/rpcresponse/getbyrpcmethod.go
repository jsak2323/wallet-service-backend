package rpcresponse

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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
		} else {
			logger.InfoLog(" - rpcresponse.GetByRpcMethodIdHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcresponse.GetByRpcMethodIdHandler, Requesting ...", req)

	if reqRpcMethodId, err = strconv.Atoi(vars["rpc_method_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcMethodIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.RpcResponses, err = s.rrsRepo.GetByRpcMethodId(reqRpcMethodId); err != nil {
		logger.ErrorLog(" -- rpcresponse.GetByRpcMethodIdHandler rrsRepo.GetByRpcMethodId Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
