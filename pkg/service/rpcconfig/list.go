package rpcconfig

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) List(ctx context.Context, page, limit int) (resp []domain.RpcConfig, err error) {
	if resp, err = s.rcRepo.GetAll(ctx, page, limit); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllRPCConfig)
		return resp, nil
	}

	return resp, nil
}
