package user

import (
	"context"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *userService) CreateUserRole(ctx context.Context, urReq userHandler.UserRoleReq) (err error) {
	if err = svc.validator.Validate(urReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = svc.urRepo.Create(ctx, urReq.UserId, urReq.RoleId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRoleUser)
		return err
	}
	return err
}
