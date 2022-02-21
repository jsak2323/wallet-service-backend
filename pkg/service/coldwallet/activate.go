package coldwallet

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *coldWalletService) ActivateColdWallet(ctx context.Context, id int) (err error) {
	if err = s.cbRepo.ToggleActive(ctx, id, true); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeactivatedColdBalance)
		return err
	}

	return nil
}
