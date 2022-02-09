package role

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
)

func (s *roleService) UpdateRole(ctx context.Context, updateReq roleHandler.UpdateReq) (err error) {
	if err = s.validator.Validate(updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.roleRepo.Update(ctx, updateReq.Role); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRole)
		return err
	}
	return nil
}
