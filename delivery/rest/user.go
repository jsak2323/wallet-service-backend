package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
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

func (re *Rest) ActivateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
		RES    StandardRes
		err    error
		ctx    = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "User successfully activated"
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

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.User
	if err = service.ActivateUser(ctx, userId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

}

func (re *Rest) DeactivateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId int
		RES    StandardRes
		err    error
		ctx    = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "User successfully deactivated"
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
			RES.Success = false
			RES.Message = ""
		}

		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.User
	if err = service.DeactivateUser(ctx, userId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

}

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

func (re *Rest) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq userHandler.CreateReq
		RES       userHandler.CreateRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Message = "User successfully created"
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

	service := re.svc.User
	id, err := service.CreateUser(ctx, createReq)
	if err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

	RES.Id = id
}

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

func (re *Rest) AddRolesHandler(w http.ResponseWriter, req *http.Request) {
	var (
		urReq userHandler.UserRoleReq
		RES   StandardRes
		err   error
		ctx   = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Message = "Role successfully added to User"
		RES.Success = true
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

	if err = json.NewDecoder(req.Body).Decode(&urReq); err != nil {

		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.User
	if err = service.CreateUserRole(ctx, urReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}

func (re *Rest) DeleteRoleHandler(w http.ResponseWriter, req *http.Request) {
	var (
		userId, roleId int
		RES            StandardRes
		err            error
		ctx            = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "Role successfully removed from User"
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

	vars := mux.Vars(req)
	if userId, err = strconv.Atoi(vars["user_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if roleId, err = strconv.Atoi(vars["role_id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.User
	if err = service.DeleteUserRole(ctx, userId, roleId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
