package wallet

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *walletService) SetPercent(ctx context.Context, res *GetBalanceRes) {
	var (
		err      error
		errField *errs.Error = nil
		hotCold  string
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if res.HotPercent, err = util.PercentCurrency(res.TotalHotCoin, res.TotalUserCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}

	if hotCold, err = util.AddCurrency(res.TotalColdCoin, res.TotalHotCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}

	if res.HotColdPercent, err = util.PercentCurrency(hotCold, res.TotalUserCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}
}
