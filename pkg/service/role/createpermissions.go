package role

import (
	"context"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *roleService) CreateRolePermission(ctx context.Context, req roleHandler.RolePermissionReq) (err error) {
	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rpRepo.Create(ctx, req.RoleId, req.PermissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRolePermission)
		return err
	}

	return nil
}
