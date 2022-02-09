package role

import (
	"context"

	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *roleService) CreateRole(ctx context.Context, createReq roleHandler.CreateReq) (id int, err error) {
	// rest to service
	// preprocessing
	if err = s.validator.Validate(createReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return id, err
	}

	if id, err = s.roleRepo.Create(ctx, createReq.Name); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRole)
		return id, err
	}
	return id, nil
}
