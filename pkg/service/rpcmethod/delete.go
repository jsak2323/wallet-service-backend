package rpcmethod

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcMethodService) Delete(ctx context.Context, id, RpcConfigId int) (err error) {

	if err = s.rcrmRepo.DeleteByRpcMethod(ctx, id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCConfigRPCMethodByRPCMethodID)
		return err
	}

	config.LoadRpcMethodByRpcConfigId(ctx, s.rmRepo, RpcConfigId)

	return nil
}
