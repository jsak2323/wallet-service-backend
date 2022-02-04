package user

import (
	"context"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) UpdateUser(ctx context.Context, updateReq userHandler.UpdateReq) (err error) {
	if !updateReq.Valid() {
		err = errs.AddTrace(errs.InvalidRequest)
		return err
	}

	if updateReq.Password != "" {
		hashPasswordByte, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGeneratePassword)
			return err
		}

		updateReq.Password = string(hashPasswordByte)
	}

	if err = s.userRepo.Update(ctx, updateReq.User); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateUser)
		return err
	}
	return err
}
