package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (re *Rest) DeleteRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId, permissionId int
		RES                  handlerRpcConfig.StandardRes
		err                  error
		ctx                  = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "RPC Method successfully removed from RPC Config"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if permissionId, err = strconv.Atoi(vars["permission_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcConfig
	if err = service.DeleteRpcMethod(ctx, roleId, permissionId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
