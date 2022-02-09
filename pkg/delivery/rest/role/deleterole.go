package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (re *Rest) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		roleId int
		RES    roleHandler.StandardRes
		err    error
		ctx    = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Role successfully deleted"
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
	if roleId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.Role
	if err = service.DeleteRole(ctx, roleId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
