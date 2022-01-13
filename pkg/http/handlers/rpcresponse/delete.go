package rpcresponse

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcresponse.DeleteHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Response successfully deleted"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, RpcMethodId)
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

	if err = s.rrsRepo.Delete(id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCResponseByID)
		return
	}
}
