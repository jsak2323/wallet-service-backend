package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RpcConfigService) DeleteRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId, permissionId int
		RES   				 StandardRes
		err   				 error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "RPC Method successfully removed from RPC Config"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		logger.ErrorLog(" - DeleteRpcMethodHandler invalid request")
		RES.Error = "Invalid request role_id"
	}

    if permissionId, err = strconv.Atoi(vars["permission_id"]); err != nil {
		logger.ErrorLog(" - DeleteRpcMethodHandler invalid request")
		RES.Error = "Invalid request permission_id"
	}

	if err = svc.rcrmRepo.Delete(roleId, permissionId); err != nil {
		logger.ErrorLog(" - AddRpcMethodsHandler svc.rcrmRepo.Delete err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}