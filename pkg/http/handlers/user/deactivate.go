package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *UserService) DeactivateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
		RES   	     StandardRes
		err   	     error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "User successfully deactivated"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    if userId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" - DeactivateUserHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.userRepo.ToggleActive(userId, false); err != nil {
		logger.ErrorLog(" - DeactivateUserHandler svc.userRepo.ToggleActive err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}