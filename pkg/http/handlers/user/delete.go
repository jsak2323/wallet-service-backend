package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *UserService) DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
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
    if userId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - DeleteUserHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.urRepo.DeleteByUserId(userId); err != nil {
		logger.ErrorLog(" - DeleteUserHandler svc.userRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if err = svc.userRepo.Delete(userId); err != nil {
		logger.ErrorLog(" - DeleteUserHandler svc.userRepo.Delete err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}