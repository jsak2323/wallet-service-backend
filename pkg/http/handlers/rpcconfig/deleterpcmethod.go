package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *RpcConfigService) DeleteRpcMethodHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId, permissionId int
		RES                  StandardRes
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

	if err = svc.rcrmRepo.Delete(ctx, roleId, permissionId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCConfigRPCMethod)
		return
	}
}
