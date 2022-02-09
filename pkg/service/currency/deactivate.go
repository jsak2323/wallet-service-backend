package currency

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) DeactivateCurrencyConfig(ctx context.Context, userId int) (err error) {
	if err = s.ccRepo.ToggleActive(ctx, userId, false); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeactivateCurrencyConfig)
		return err
	}
	return nil
}
