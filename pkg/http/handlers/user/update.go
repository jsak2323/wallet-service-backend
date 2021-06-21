package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		updateReq   UpdateReq
		RES         StandardRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		logger.ErrorLog(" - UpdateUserHandler json.NewDecoder err: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if !updateReq.valid() {
		logger.ErrorLog(" - UpdateUserHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	if updateReq.Password != "" {
		hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.ErrorLog(" - UpdateUserHandler bcrypt.GenerateFromPassword err: " + err.Error())
			RES.Error = err.Error()
			return
		}

		updateReq.Password = string(hashPasswordByte)
	}

	if err = svc.userRepo.Update(updateReq.User); err != nil {
		logger.ErrorLog(" - UpdateUserHandler svc.roleRepo.Update err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}