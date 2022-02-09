package rpcrequest

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcRequestService) Delete(ctx context.Context, id, RpcMethodId int) (err error) {

	if err = s.rrqRepo.Delete(id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteRPCRequestByID)
		return err
	}

	config.LoadRpcRequestByRpcMethodId(s.rrqRepo, RpcMethodId)

	return nil
}
