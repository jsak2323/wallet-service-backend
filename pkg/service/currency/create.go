package currency

import (
	"context"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) CreateCurrencyConfig(ctx context.Context, req domain.CurrencyConfig) (err error) {
	if err = s.validator.Validate(req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.ccRepo.Create(ctx, req); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateCurrencyConfig)
		return err
	}
	return nil
}
