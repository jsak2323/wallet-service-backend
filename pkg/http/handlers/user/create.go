package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)


func (svc *UserService) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq CreateReq
		RES       CreateRes
		err       error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Message = "User successfully created"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		logger.ErrorLog(" - CreateUserHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if !createReq.valid() {
		logger.ErrorLog(" - CreateUserHandler invalid request")
		RES.Error = "Invalid request"
		return
	}

	hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLog(" - CreateUserHandler bcrypt.GenerateFromPassword err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	createReq.Password = string(hashPasswordByte)

	if RES.Id, err = svc.userRepo.Create(createReq.User); err != nil {
		logger.ErrorLog(" - CreateUserHandler svc.userRepo.Create err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}