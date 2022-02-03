package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
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

func (re *Rest) CreatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq CreateReq
		RES       CreateRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Message = "Permission successfully created"
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

	service := re.svc.Permission
	if RES.Id, err = service.CreatePermission(ctx, createReq.Name); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

}

func (re *Rest) UpdatePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq domain.Permission
		RES       StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Permission successfully updated"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
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

	service := re.svc.Permission
	err = service.UpdatePermission(ctx, updateReq)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}

func (re *Rest) DeletePermissionHandler(w http.ResponseWriter, req *http.Request) {
	var (
		permissionId int
		RES          StandardRes
		err          error
		ctx          = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Permission successfully deleted"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			RES.Success = false
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if permissionId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.Permission
	if err = service.DeletePermission(ctx, permissionId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
