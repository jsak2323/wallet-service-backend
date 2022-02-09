package rpcconfig

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) GetById(ctx context.Context, idRpcConfig int) (res domain.RpcConfig, err error) {

	if res, err = s.rcRepo.GetById(ctx, idRpcConfig); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCConfigByID)
		return res, err
	}
	return res, nil
}
