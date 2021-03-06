package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId, permissionId int
		RES                  StandardRes
		err                  error
		ctx                  = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Permission successfully removed from Role"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
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

	if err = svc.rpRepo.Delete(roleId, permissionId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermission)
		return
	}
}
