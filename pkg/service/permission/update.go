package permission

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *permissionService) UpdatePermission(ctx context.Context, updateReq domain.Permission) (err error) {
	if err = svc.validator.Validate(updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = svc.permissionRepo.Update(ctx, updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdatePermission)
		return err
	}
	return nil
}
