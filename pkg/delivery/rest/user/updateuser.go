package user

import (
	"encoding/json"
	"net/http"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq userHandler.UpdateReq
		RES       StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "User successfully updated"

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

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.User
	if err = service.UpdateUser(ctx, updateReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
