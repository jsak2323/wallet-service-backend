package wallet

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceHandlerResponseMap map[int]GetBalanceRes

func (s *walletService) SetHotLimits(ctx context.Context, res *GetBalanceRes) {
	var (
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if res.HotLimits, err = s.hotLimitRepo.GetBySymbol(res.CurrencyConfig.Symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotLimits)
	}
}
