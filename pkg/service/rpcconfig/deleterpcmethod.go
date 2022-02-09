package rpcconfig

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) DeleteRpcMethod(ctx context.Context, roleId, permissionId int) (err error) {

	if err = s.rcrmRepo.Delete(ctx, roleId, permissionId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCConfigRPCMethod)
		return err
	}

	return nil
}
