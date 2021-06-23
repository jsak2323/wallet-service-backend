package role

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) CreateRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq   CreateReq
		RES         CreateRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Message = "Role successfully created"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		logger.ErrorLog(" - CreateRoleHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if !createReq.valid() {
		logger.ErrorLog(" - CreateRoleHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if RES.Id, err = svc.roleRepo.Create(createReq.Name); err != nil {
		logger.ErrorLog(" - CreateRoleHandler svc.roleRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}