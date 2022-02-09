package rpcmethod

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcMethodService) List(ctx context.Context, page, limit int) (resp []domain.RpcMethod, err error) {

	if resp, err = s.rmRepo.GetAll(ctx, page, limit); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllRPCMethod)
		return resp, err
	}

	return resp, nil
}
