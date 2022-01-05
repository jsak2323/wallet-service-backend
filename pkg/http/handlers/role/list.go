package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) ListRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	roles, err := svc.roleRepo.GetAll(page, limit)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedGetAllRole.Title})
		return
	}

	for i, role := range roles {
		roles[i].Permissions, err = svc.permissionRepo.GetByRoleId(role.Id)
		if err != nil {
			RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedGetPermissionByRole.Title})
			return
		}
	}

	RES.Roles = roles
}
