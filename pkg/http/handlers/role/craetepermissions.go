package role

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpReq   RolePermissionReq
		RES   	StandardRes
		err   	error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Permission successfully added to Role"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&rpReq); err != nil {
		logger.ErrorLog(" - CreatePermissionHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if !rpReq.valid() {
		logger.ErrorLog(" - CreatePermissionHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if err = svc.rpRepo.Create(rpReq.RoleId, rpReq.PermissionId); err != nil {
		logger.ErrorLog(" - CreatePermissionHandler svc.rpRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}