package role

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *RoleService) CreateRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq CreateReq
		RES       CreateRes
		err       error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error))
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

	if !createReq.valid() {
		RES.Error = errs.AddTrace(errs.InvalidRequest)
		return
	}

	if RES.Id, err = svc.roleRepo.Create(createReq.Name); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRole)
		return
	}
}
