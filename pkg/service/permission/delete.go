package permission

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *permissionService) DeletePermission(ctx context.Context, permissionId int) (err error) {
	if err = svc.rpRepo.DeleteByPermissionId(ctx, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermissionByPermissionID)
		return err
	}

	if err = svc.permissionRepo.Delete(ctx, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeletePermissionByID)
		return err
	}
	return nil
}
