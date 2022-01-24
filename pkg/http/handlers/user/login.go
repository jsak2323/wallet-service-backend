package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) LoginHandler(w http.ResponseWriter, req *http.Request) {
	var (
		loginReq LoginReq
		RES      LoginRes

		user user.User
		td   jwt.TokenDetails
		err  error
		ctx  = req.Context()
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

	if err = svc.validator.Validate(loginReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if user, err = svc.userRepo.GetByUsername(loginReq.Username); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.UsernameNotFound)
		return
	}

	if user.RoleNames, err = svc.roleRepo.GetNamesByUserId(user.Id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRolesByUserID)
		return
	}

	if user.PermissionNames, err = svc.permissionRepo.GetNamesByUserId(user.Id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.Permissions)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.IncorrectPassword)
		return
	}

	if td, err = jwt.CreateToken(user); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateToken)
		return
	}

	RES.AccessToken = td.AccessToken
	RES.RefreshToken = td.RefreshToken
}
