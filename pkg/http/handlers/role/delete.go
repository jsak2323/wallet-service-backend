package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *RoleService) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId int
		RES   	     StandardRes
		err   	     error
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
    if roleId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - DeleteRoleHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.urRepo.DeleteByRoleId(roleId); err != nil {
		logger.ErrorLog(" - DeleteRoleHandler svc.roleRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = svc.rpRepo.DeleteByRoleId(roleId); err != nil {
		logger.ErrorLog(" - DeleteRoleHandler svc.roleRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = svc.roleRepo.Delete(roleId); err != nil {
		logger.ErrorLog(" - DeleteRoleHandler svc.roleRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}