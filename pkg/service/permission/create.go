package permission

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *permissionService) CreatePermission(ctx context.Context, name string) (id int, err error) {
	// if err = s.validator.Validate(createReq); err != nil {
	// 	RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
	// 	return id, err
	// }

	if name == "" {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return id, err
	}

	if id, err = s.permissionRepo.Create(ctx, name); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreatePermission)
		return id, err
	}
	return id, err
}
