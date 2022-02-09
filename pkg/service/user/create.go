package user

import (
	"context"

	domainUser "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, createReq domainUser.CreateReq) (id int, err error) {
	if err = s.validator.Validate(createReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return id, err
	}

	hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGeneratePassword)
		return id, err
	}

	createReq.Password = string(hashPasswordByte)

	if id, err = s.userRepo.Create(ctx, createReq.User); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateUser)
		return id, err
	}
	return id, nil
}
