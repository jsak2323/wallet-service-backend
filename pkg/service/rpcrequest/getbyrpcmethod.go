package rpcrequest

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcRequestService) GetByRpcMethodId(ctx context.Context, reqRpcMethodId int) (resp []domain.RpcRequest, err error) {

	if resp, err = s.rrqRepo.GetByRpcMethodId(reqRpcMethodId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCRequestByRPCMethodID)
		return resp, err
	}

	return resp, nil
}
