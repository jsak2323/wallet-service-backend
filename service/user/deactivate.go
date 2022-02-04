package user

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *userService) DeactivateUser(ctx context.Context, userId int) (err error) {
	if err = s.userRepo.ToggleActive(ctx, userId, false); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedActivateUser)
		return err
	}
	return nil
}
