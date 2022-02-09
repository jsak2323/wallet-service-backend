package rpcmethod

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcMethodService) Create(ctx context.Context, req domain.RpcMethod) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if req.Id, err = s.rmRepo.Create(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCMethod)
		return err
	}

	if err = s.rcrmRepo.Create(ctx, req.RpcConfigId, req.Id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfigRPCMethod)
		return err
	}

	config.LoadRpcMethodByRpcConfigId(ctx, s.rmRepo, req.RpcConfigId)

	return nil
}
