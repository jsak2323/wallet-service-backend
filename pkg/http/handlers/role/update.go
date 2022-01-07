package role

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) UpdateRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq UpdateReq
		RES       StandardRes
		err       error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Role successfully updated"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error))
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

	if !updateReq.valid() {
		RES.Error = errs.AddTrace(errs.InvalidRequest)
		return
	}

	if err = svc.roleRepo.Update(updateReq.Role); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRole)
		return
	}
}
