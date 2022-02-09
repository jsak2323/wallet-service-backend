package rpcconfig

import (
	"context"

	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) CreateRpcMethod(ctx context.Context, req handlerRpcConfig.RpcConfigRpcMethodReq) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rcrmRepo.Create(ctx, req.RpcConfigId, req.RpcMethodId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfigRPCMethod)
		return err
	}

	return nil
}
