package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (svc *RoleService) ListRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES         ListRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			RES.Error = err.Error()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])
	
	roles, err := svc.roleRepo.GetAll(page, limit)
	if err != nil {
		return
	}

	for i, user := range roles {
		roles[i].Permissions, err = svc.permissionRepo.GetByRoleId(user.Id)
		if err != nil {
			return
		}
	}

	RES.Roles = roles
}