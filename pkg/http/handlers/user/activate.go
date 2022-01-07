package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *UserService) ActivateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
		RES    StandardRes
		err    error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "User successfully activated"
		if RES.Error != nil {
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
	if userId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = svc.userRepo.ToggleActive(userId, true); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedActivateUser)
		return
	}
}
