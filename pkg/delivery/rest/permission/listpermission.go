package permission

import (
	"encoding/json"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (re *Rest) ListPermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
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

	service := re.svc.Permission
	RES.Permissions, err = service.ListPermissions(ctx, page, limit)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllPermission)
		return
	}
}
