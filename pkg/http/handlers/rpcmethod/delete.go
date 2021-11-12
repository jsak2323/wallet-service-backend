package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES             StandardRes
		err             error
		id, RpcConfigId int
	)

	vars := mux.Vars(req)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcmethod.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully deleted"

			config.LoadRpcMethodByRpcConfigId(s.rmRepo, RpcConfigId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.DeleteHandler, Requesting ...", req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcConfigIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RpcConfigId, err = strconv.Atoi(vars["rpc_config_id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByRpcConfigIdHandler strconv.Atoi(" + vars["rpc_config_id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = s.rcrmRepo.DeleteByRpcMethod(id); err != nil {
		logger.ErrorLog(" -- rpcmethod.DeleteHandler rcrmRepo.DeleteByRpcMethod Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = s.rmRepo.Delete(id); err != nil {
		logger.ErrorLog(" -- rpcmethod.DeleteHandler rmRepo.Delete Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
