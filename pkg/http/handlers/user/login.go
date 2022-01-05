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
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}

	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if user, err = svc.userRepo.GetByUsername(loginReq.Username); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.UsernameNotFound.Title})
		return
	}

	if user.RoleNames, err = svc.roleRepo.GetNamesByUserId(user.Id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedGetRolesByUserID.Title})
		return
	}

	if user.PermissionNames, err = svc.permissionRepo.GetNamesByUserId(user.Id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.Permissions.Title})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.IncorrectPassword.Title})
		return
	}

	if td, err = jwt.CreateToken(user); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedCreateToken.Title})
		return
	}

	RES.AccessToken = td.AccessToken
	RES.RefreshToken = td.RefreshToken
}
