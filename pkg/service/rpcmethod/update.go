package rpcmethod

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcMethodService) Update(ctx context.Context, req domain.UpdateRpcMethod) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rmRepo.Update(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCMethod)
		return err
	}

	config.LoadRpcMethodByRpcConfigId(ctx, s.rmRepo, req.RpcConfigId)
	return nil
}
