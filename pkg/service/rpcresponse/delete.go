package rpcresponse

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcResponseService) Delete(ctx context.Context, id, rpcMethodId int) (err error) {
	if err = s.rpcresponseRepo.Delete(ctx, id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCResponseByID)
		return err
	}

	config.LoadRpcResponseByRpcMethodId(ctx, s.rpcresponseRepo, rpcMethodId)

	return nil
}
