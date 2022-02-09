package user

import (
	"encoding/json"
	"net/http"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) AddUserRolesHandler(w http.ResponseWriter, req *http.Request) {
	var (
		urReq userHandler.UserRoleReq
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

	service := re.svc.User
	if err = service.CreateUserRole(ctx, urReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
