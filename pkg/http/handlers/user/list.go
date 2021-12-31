package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) ListUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES      ListRes
		err      error
		errTitle string
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

	users, err := svc.userRepo.GetAll(page, limit)
	if err != nil {
		errTitle = errs.FailedGetAllUser.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	for i, user := range users {
		users[i].Roles, err = svc.roleRepo.GetByUserId(user.Id)
		if err != nil {
			errTitle = errs.FailedGetRoleByUserId.Title
			RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
			return
		}
	}

	RES.Users = users
}
