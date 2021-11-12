package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) GetByRpcConfigIdHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error
	var reqRpcConfigId int

	vars := mux.Vars(req)
	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcmethod.ListHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcmethod.ListHandler, Requesting ...", req)

	if reqRpcConfigId, err = strconv.Atoi(vars["rpc_config_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.RpcMethods, err = s.rmRepo.GetByRpcConfigId(reqRpcConfigId); err != nil {
		logger.ErrorLog(" -- rpcmethod.ListHandler s.rmRepo.GetByRpcConfigId Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
