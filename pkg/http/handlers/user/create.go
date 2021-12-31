package user

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (svc *UserService) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	var (
		createReq CreateReq
		RES       CreateRes
		err       error
		errTitle  string
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		RES.Message = "User successfully created"
		if err != nil {
			resStatus = http.StatusInternalServerError
			RES.Message = ""
			logger.ErrorLog(errs.Logged(RES.Error))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		errTitle = errs.ErrorUnmarshalBodyRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	if !createReq.valid() {
		err = errs.InvalidRequest
		errTitle = errs.InvalidRequest.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		errTitle = errs.FailedGeneratePassword.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}

	createReq.Password = string(hashPasswordByte)

	if RES.Id, err = svc.userRepo.Create(createReq.User); err != nil {
		log.Println(err)
		errTitle = errs.FailedCreateUser.Title
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errTitle})
		return
	}
}
