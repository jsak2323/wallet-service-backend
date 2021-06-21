package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *RoleService) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId, permissionId int
		RES   				 StandardRes
		err   				 error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler invalid request")
		RES.Error = "Invalid request role_id"
	}

    if permissionId, err = strconv.Atoi(vars["permission_id"]); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler invalid request")
		RES.Error = "Invalid request permission_id"
	}

	if err = svc.rpRepo.Delete(roleId, permissionId); err != nil {
		logger.ErrorLog(" - AddPermissionsHandler svc.rpRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}