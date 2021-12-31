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
		td       jwt.TokenDetails
		err      error
		errTitle string
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
		errTitle = errs.ErrorUnmarshalBodyRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if user, err = svc.userRepo.GetByUsername(loginReq.Username); err != nil {
		errTitle = errs.UsernameNotFound.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if user.RoleNames, err = svc.roleRepo.GetNamesByUserId(user.Id); err != nil {
		errTitle = errs.RolesNotFound.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if user.PermissionNames, err = svc.permissionRepo.GetNamesByUserId(user.Id); err != nil {
		errTitle = errs.Permissions.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		errTitle = errs.IncorrectPassword.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if td, err = jwt.CreateToken(user); err != nil {
		errTitle = errs.FailedCreateToken.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	RES.AccessToken = td.AccessToken
	RES.RefreshToken = td.RefreshToken
}
