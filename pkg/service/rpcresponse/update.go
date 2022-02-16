package rpcresponse

import (
	"context"
	"encoding/json"
	"errors"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcResponseService) Update(ctx context.Context, req domain.RpcResponse) (err error) {
	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	err = json.Unmarshal([]byte(req.JsonFieldsStr), req.JsonFields)
	if err != nil && req.JsonFieldsStr != "" {
		return errs.AddTrace(errors.New("JSON data at JSON Field"))
	}

	if err = s.rpcresponseRepo.Update(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCResponse)
		return err
	}
	return err
}
