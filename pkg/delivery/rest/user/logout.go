package user

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES   StandardRes
		valid bool
		err   error
		ctx   = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
			RES.Success = false
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if _, valid, err = jwt.ParseFromRequest(req); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateToken)
		return
	}

	if !valid {
		RES.Error = errs.AddTrace(errs.InvalidToken)
		return
	}
}
