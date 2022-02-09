package permission

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *permissionService) UpdatePermission(ctx context.Context, updateReq domain.Permission) (err error) {
	if err = s.validator.Validate(updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.permissionRepo.Update(ctx, updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdatePermission)
		return err
	}
	return nil
}
