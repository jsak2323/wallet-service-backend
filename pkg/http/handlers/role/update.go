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
		if err != nil {
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
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if !updateReq.valid() {
		err = errs.InvalidRequest
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = svc.roleRepo.Update(updateReq.Role); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateRole.Title})
		return
	}
}
