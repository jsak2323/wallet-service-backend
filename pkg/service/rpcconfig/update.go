package rpcconfig

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) Update(ctx context.Context, req domain.UpdateRpcConfig) (err error) {

	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.rcRepo.Update(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCConfig)
		return err
	}

	config.LoadCurrencyConfigs(ctx)

	return nil
}
