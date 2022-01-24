package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) GetByRpcConfigIdHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES            ListRes
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

	if RES.RpcMethods, err = s.rmRepo.GetByRpcConfigId(reqRpcConfigId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCMethodByConfigID)
		return
	}
}
