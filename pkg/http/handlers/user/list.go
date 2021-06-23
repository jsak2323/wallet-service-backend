package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) ListUserHandler(w http.ResponseWriter, req *http.Request) {
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
	
	users, err := svc.userRepo.GetAll(page, limit)
	if err != nil {
		logger.ErrorLog(" - ListUserHandler svc.userRepo.GetAll err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	for i, user := range users {
		users[i].Roles, err = svc.roleRepo.GetByUserId(user.Id)
		if err != nil {
			logger.ErrorLog(" - ListUserHandler svc.roleRepo.GetByUserId err: " + err.Error())
			RES.Error = errInternalServer
			return
		}
	}

	RES.Users = users
}