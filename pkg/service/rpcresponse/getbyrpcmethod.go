package rpcresponse

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcResponseService) GetByRpcMethod(ctx context.Context, reqRpcMethodId int) (resp []rpcresponse.RpcResponse, err error) {
	if resp, err = s.rpcresponseRepo.GetByRpcMethodId(ctx, reqRpcMethodId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCResponseByRPCMethodID)
		return resp, err
	}

	return resp, nil
}
