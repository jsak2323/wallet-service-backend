package permission

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *PermissionService) UpdatePermissionHandler(w http.ResponseWriter, req *http.Request) {
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
		logger.ErrorLog(" - UpdatePermissionHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !updateReq.valid() {
		logger.ErrorLog(" - UpdatePermissionHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if err = svc.permissionRepo.Update(updateReq.Permission); err != nil {
		logger.ErrorLog(" - UpdatePermissionHandler svc.permissionRepo.Update err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}