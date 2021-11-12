package rpcresponse

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcResponseService) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES             StandardRes
		err             error
		id, RpcMethodId int
	)

	vars := mux.Vars(req)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcresponse.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Response successfully deleted"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.DeleteHandler, Responseing ...", req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcMethodIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RpcMethodId, err = strconv.Atoi(vars["rpc_method_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcMethodIdHandler strconv.Atoi(" + vars["rpc_method_id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = s.rrsRepo.Delete(id); err != nil {
		logger.ErrorLog(" -- rpcresponse.DeleteHandler rrsRepo.Delete Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
