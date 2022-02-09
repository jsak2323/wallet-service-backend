package user

import (
	"encoding/json"
	"net/http"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) LoginHandler(w http.ResponseWriter, req *http.Request) {
	var (
		loginReq userHandler.LoginReq
		RES      userHandler.LoginRes
		err      error
		ctx      = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}

	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.User

	res, err := service.Login(ctx, loginReq)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

	RES = res
}
