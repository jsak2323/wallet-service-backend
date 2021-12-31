package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq UpdateReq
		RES       StandardRes
		err       error
		errTitle  string
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Success = true
		RES.Message = "User successfully updated"

		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
			RES.Success = false
			RES.Message = ""
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}

	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		errTitle = errs.ErrorUnmarshalBodyRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if !updateReq.valid() {
		errTitle = errs.InvalidRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if updateReq.Password != "" {
		hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			errTitle = errs.FailedGeneratePassword.Title
			RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
			return
		}

		updateReq.Password = string(hashPasswordByte)
	}

	if err = svc.userRepo.Update(updateReq.User); err != nil {
		errTitle = errs.FailedUpdateUser.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}
}
