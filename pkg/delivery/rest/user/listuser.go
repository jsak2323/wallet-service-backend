package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (re *Rest) ListUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES userHandler.ListRes
		err error
		ctx = req.Context()
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

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	service := re.svc.User
	users, err := service.ListUser(ctx, page, limit)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

	RES.Users = users
}
