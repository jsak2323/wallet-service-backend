package user

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) AddRolesHandler(w http.ResponseWriter, req *http.Request) {
	var (
		urReq UserRoleReq
		RES   StandardRes
		err   error
		ctx   = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Message = "Role successfully added to User"
		RES.Success = true
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
			RES.Success = false
			RES.Message = ""
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&urReq); err != nil {

		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = svc.validator.Validate(urReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = svc.urRepo.Create(ctx, urReq.UserId, urReq.RoleId); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRoleUser)
		return
	}
}
