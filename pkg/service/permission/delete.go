package permission

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *permissionService) DeletePermission(ctx context.Context, permissionId int) (err error) {
	if err = s.rpRepo.DeleteByPermissionId(ctx, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermissionByPermissionID)
		return err
	}

	if err = s.permissionRepo.Delete(ctx, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeletePermissionByID)
		return err
	}
	return nil
}
