package rpcresponse

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcResponseService) Create(ctx context.Context, req domain.CreateRpcResponse) (err error) {
	if err = s.validator.Validate(req); err != nil && req.JsonFieldsStr != "" {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	err = json.Unmarshal([]byte(req.JsonFieldsStr), req.JsonFields)
	if err != nil && req.JsonFieldsStr != "" {
		err = errs.AddTrace(errors.New("JSON data at JSON Field"))
		return err
	}

	if err = s.rpcresponseRepo.Create(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCResponse)
		return err
	}

	config.LoadRpcResponseByRpcMethodId(ctx, s.rpcresponseRepo, req.RpcMethodId)

	return nil
}
