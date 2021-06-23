package permission

import (
	"encoding/json"
	"net/http"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *PermissionService) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
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
			RES.Message = "Permission successfully created"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		logger.ErrorLog(" - CreatePermissionHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if !createReq.valid() {
		logger.ErrorLog(" - CreatePermissionHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if RES.Id, err = svc.permissionRepo.Create(createReq.Name); err != nil {
		logger.ErrorLog(" - CreatePermissionHandler svc.permissionRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}