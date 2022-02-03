package permission

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (svc *permissionService) CreatePermission(ctx context.Context, name string) (id int, err error) {
	// if err = svc.validator.Validate(createReq); err != nil {
	// 	RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
	// 	return id, err
	// }

	if name == "" {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return id, err
	}

	if id, err = svc.permissionRepo.Create(ctx, name); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreatePermission)
		return id, err
	}
	return id, err
}
