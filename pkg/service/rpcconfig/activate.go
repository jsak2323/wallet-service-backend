package rpcconfig

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *rpcConfigService) Activate(ctx context.Context, id int) (err error) {

	if err = s.rcRepo.ToggleActive(ctx, id, true); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedActivateRPCConfig)
		return err
	}

	config.LoadAppConfig()

	return nil
}
