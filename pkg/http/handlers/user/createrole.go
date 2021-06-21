package user

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) AddRolesHandler(w http.ResponseWriter, req *http.Request) {
	var (
		urReq   UserRoleReq
		RES   	StandardRes
		err   	error
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

	if err = json.NewDecoder(req.Body).Decode(&urReq); err != nil {
		logger.ErrorLog(" - AddRolesHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !urReq.valid() {
		logger.ErrorLog(" - AddRolesHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if err = svc.urRepo.Create(urReq.UserId, urReq.RoleId); err != nil {
		logger.ErrorLog(" - AddRolesHandler svc.urRepo.Create err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}