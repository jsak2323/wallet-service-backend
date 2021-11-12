package rpcrequest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcRequestService) DeleteHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- rpcrequest.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully deleted"

			config.LoadRpcRequestByRpcMethodId(s.rrqRepo, RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcrequest.DeleteHandler, Requesting ...", req)

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

	if err = s.rrqRepo.Delete(id); err != nil {
		logger.ErrorLog(" -- rpcrequest.DeleteHandler rrqRepo.Delete Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
