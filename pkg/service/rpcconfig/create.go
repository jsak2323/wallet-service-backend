package rpcconfig

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) Create(ctx context.Context, req domain.RpcConfig) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rcRepo.Create(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfig)
		return err
	}

	config.LoadCurrencyConfigs(ctx)
	return nil
}
