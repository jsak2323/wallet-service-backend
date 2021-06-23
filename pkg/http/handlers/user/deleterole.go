package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *UserService) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId, roleId int
		RES   				 StandardRes
		err   				 error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Role successfully removed from User"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    if userId, err = strconv.Atoi(vars["user_id"]); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler invalid request")
		RES.Error = "Invalid request user_id"
	}

    if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		logger.ErrorLog(" - DeletePermissionHandler invalid request")
		RES.Error = "Invalid request role_id"
	}

	if err = svc.urRepo.Delete(userId, roleId); err != nil {
		logger.ErrorLog(" - AddPermissionsHandler svc.userRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}