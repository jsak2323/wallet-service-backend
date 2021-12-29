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

		user     user.User
		err      error
		errTitle string
		td       jwt.TokenDetails
		funcName = errs.GetFunctionName(svc.LoginHandler)
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			RES.Error = errs.AssignErr(errs.AddTrace(err, funcName), errs.Error{Title: errTitle})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		logger.ErrorLog(" - LoginHandler json.NewDecoder err: " + err.Error())
		errTitle = errs.ErrorUnmarshalBodyRequest.Title
		return
	}

	if user, err = svc.userRepo.GetByUsername(loginReq.Username); err != nil {
		logger.ErrorLog(" - LoginHandler svc.userRepo.GetByUsername err: " + err.Error())
		errTitle = errs.UsernameNotFound.Title
		return
	}

	if user.RoleNames, err = svc.roleRepo.GetNamesByUserId(user.Id); err != nil {
		logger.ErrorLog(" - LoginHandler svc.roleRepo.GetNamesByUserId err: " + err.Error())
		errTitle = errs.RolesNotFound.Title
		return
	}

	if user.PermissionNames, err = svc.permissionRepo.GetNamesByUserId(user.Id); err != nil {
		logger.ErrorLog(" - LoginHandler svc.permissionRepo.GetNamesByUserId err: " + err.Error())
		errTitle = errs.Permissions.Title
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		logger.ErrorLog(" - LoginHandler bcrypt.CompareHashAndPassword err: " + err.Error())
		errTitle = errs.IncorrectPassword.Title
		return
	}

	if td, err = jwt.CreateToken(user); err != nil {
		logger.ErrorLog(" - LoginHandler CreateToken err: " + err.Error())
		errTitle = errs.FailedCreateToken.Title
		return
	}

	RES.AccessToken = td.AccessToken
	RES.RefreshToken = td.RefreshToken
}
