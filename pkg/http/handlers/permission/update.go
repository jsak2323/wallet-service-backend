package permission

import (
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *PermissionService) UpdatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq UpdateReq
		RES       StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Permission successfully updated"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	err = UpdatePermission(err, svc, updateReq, RES, ctx)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}

func UpdatePermission(err error, svc *PermissionService, updateReq UpdateReq, RES StandardRes, ctx context.Context) error {
	if err = svc.validator.Validate(updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return RES.Error
	}

	if err = svc.permissionRepo.Update(ctx, updateReq.Permission); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdatePermission)
		return RES.Error
	}
	return nil
}
