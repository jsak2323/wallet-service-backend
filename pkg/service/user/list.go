package user

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *userService) ListUser(ctx context.Context, page int, limit int) (res []user.User, err error) {
	users, err := svc.userRepo.GetAll(page, limit)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllUser)
		return res, err
	}

	for i, user := range users {
		users[i].Roles, err = svc.roleRepo.GetByUserId(user.Id)
		if err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRoleByUserID)
			return res, err
		}
	}
	return users, nil
}
