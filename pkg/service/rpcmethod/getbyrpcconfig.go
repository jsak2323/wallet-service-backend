package rpcmethod

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcMethodService) GetByRpcConfigId(ctx context.Context, reqRpcConfigId int) (resp []domain.RpcMethod, err error) {

	if resp, err = s.rmRepo.GetByRpcConfigId(ctx, reqRpcConfigId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCMethodByConfigID)
		return resp, err
	}

	return resp, nil
}
