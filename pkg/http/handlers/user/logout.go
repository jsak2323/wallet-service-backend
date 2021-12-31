package user

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES StandardRes

		valid    bool
		err      error
		errTitle string
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
			RES.Success = false
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if _, valid, err = jwt.ParseFromRequest(req); err != nil {
		errTitle = errs.FailedCreateToken.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if !valid {
		errTitle = errs.InvalidToken.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}
}
