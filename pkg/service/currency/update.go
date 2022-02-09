package currency

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) UpdateCurrencyConfig(ctx context.Context, currencyConfig domain.UpdateCurrencyConfig) (err error) {
	if err = s.validator.Validate(currencyConfig); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.ccRepo.Update(ctx, currencyConfig); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateCurrencyConfig)
		return err
	}
	return nil
}
