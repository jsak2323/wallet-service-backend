package role

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *roleService) DeleteRolePermission(ctx context.Context, roleId int, permissionId int) (err error) {
	if err = s.rpRepo.Delete(roleId, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermission)
		return err
	}
	return nil
}
