package permission

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *PermissionService) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq CreateReq
		RES       CreateRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Message = "Permission successfully created"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = svc.validator.Validate(createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if RES.Id, err = svc.permissionRepo.Create(ctx, createReq.Name); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreatePermission)
		return
	}
}
