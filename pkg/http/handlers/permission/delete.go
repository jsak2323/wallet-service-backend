package permission

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *PermissionService) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		permissionId int
		RES   	     StandardRes
		err   	     error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Permission successfully deleted"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    if permissionId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.rpRepo.DeleteByPermissionId(permissionId); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler svc.rpRepo.DeleteByPermissionId err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = svc.permissionRepo.Delete(permissionId); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler svc.rpRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}