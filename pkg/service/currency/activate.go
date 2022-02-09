package currency

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) ActivateCurrency(ctx context.Context, id int) (err error) {
	if err = s.ccRepo.ToggleActive(ctx, id, true); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedActivateCurrencyConfig)
		return err
	}
	return nil
}
