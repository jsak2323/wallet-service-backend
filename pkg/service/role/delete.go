package role

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *roleService) DeleteRole(ctx context.Context, roleId int) (err error) {
	if err = s.urRepo.DeleteByRoleId(ctx, roleId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRoleUser)
		return err
	}

	if err = s.rpRepo.DeleteByRoleId(ctx, roleId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRolePermission)
		return err
	}

	if err = s.roleRepo.Delete(ctx, roleId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRole)
		return err
	}

	return nil
}
