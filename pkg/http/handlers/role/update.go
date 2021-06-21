package role

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) UpdateRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq   UpdateReq
		RES         StandardRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		logger.ErrorLog(" - UpdateRoleHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !updateReq.valid() {
		logger.ErrorLog(" - UpdateRoleHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if err = svc.roleRepo.Update(updateReq.Role); err != nil {
		logger.ErrorLog(" - UpdateRoleHandler svc.roleRepo.Update err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}