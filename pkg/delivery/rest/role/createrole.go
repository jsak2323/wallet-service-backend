package role

import (
	"encoding/json"
	"net/http"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) CreateRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq roleHandler.CreateReq
		RES       roleHandler.CreateRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {
		// transform
		resStatus := http.StatusOK
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

	service := re.svc.Role
	if RES.Id, err = service.CreateRole(ctx, createReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
