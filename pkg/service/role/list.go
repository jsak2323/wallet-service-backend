package role

import (
	"context"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *roleService) GetListRoleDetails(ctx context.Context, page int, limit int) (res roleHandler.ListRes, err error) {
	roles, err := s.roleRepo.GetAll(page, limit)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllRole)
		return res, err
	}

	for i, role := range roles {
		roles[i].Permissions, err = s.permissionRepo.GetByRoleId(ctx, role.Id)
		if err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetPermissionByRole)
			return res, err
		}
	}

	res.Roles = roles
	return res, err
}
