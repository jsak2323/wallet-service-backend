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

	if !updateReq.valid() {
		RES.Error = errs.AddTrace(errs.InvalidRequest)
		return
	}

	if updateReq.Password != "" {
		hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGeneratePassword)
			return
		}

		updateReq.Password = string(hashPasswordByte)
	}

	if err = svc.userRepo.Update(updateReq.User); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateUser)
		return
	}
}
