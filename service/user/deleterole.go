package user

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *userService) DeleteUserRole(ctx context.Context, userId int, roleId int) (err error) {
	if err = s.urRepo.Delete(ctx, userId, roleId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRoleUser)
		return err
	}
	return nil
}
