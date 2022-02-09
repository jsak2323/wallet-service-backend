package permission

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *permissionService) ListPermissions(ctx context.Context, page int, limit int) (res []permission.Permission, err error) {
	res, err = s.permissionRepo.GetAll(ctx, page, limit)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllPermission)
		return res, err
	}
	return res, err
}
