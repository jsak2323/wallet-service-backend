package user

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)


func (svc *UserService) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES StandardRes

		valid  bool
		err    error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if _, valid, err = jwt.ParseFromRequest(req); err != nil {
		logger.ErrorLog(" - LogoutHandler jwt.ParseFromRequest err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !valid {
		logger.ErrorLog(" - LogoutHandler invalid token")
		RES.Error = err.Error()
		return
	}
}