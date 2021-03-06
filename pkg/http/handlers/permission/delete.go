package permission

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *PermissionService) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		permissionId int
		RES          StandardRes
		err          error
		ctx          = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Permission successfully deleted"
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
	if permissionId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = svc.rpRepo.DeleteByPermissionId(ctx, permissionId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermissionByPermissionID)
		return
	}

	if err = svc.permissionRepo.Delete(ctx, permissionId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedDeletePermissionByID)
		return
	}
}
