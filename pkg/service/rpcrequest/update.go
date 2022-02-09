package rpcrequest

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcRequestService) Update(ctx context.Context, req domain.UpdateRpcRequest) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rrqRepo.Update(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCRequest)
		return err
	}

	config.LoadRpcRequestByRpcMethodId(s.rrqRepo, req.RpcMethodId)

	return nil
}
