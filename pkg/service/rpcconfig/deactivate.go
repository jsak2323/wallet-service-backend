package rpcconfig

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) Deactivate(ctx context.Context, idRpcConfig int) (err error) {

	if err = s.rcRepo.ToggleActive(ctx, idRpcConfig, false); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeactivateRPCConfig)
		return err
	}

	config.LoadAppConfig()

	return nil
}
