package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *UserService) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId, roleId int
		RES            StandardRes
		err            error
		errTitle       string
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Role successfully removed from User"
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
			RES.Success = false
			RES.Message = ""
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["user_id"]); err != nil {
		errTitle = errs.InvalidRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
	}

	if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		errTitle = errs.InvalidRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
	}

	if err = svc.urRepo.Delete(userId, roleId); err != nil {
		errTitle = errs.FailedDeleteRoleUser.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}
}
