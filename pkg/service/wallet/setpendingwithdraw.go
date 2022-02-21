package wallet

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *walletService) SetPendingWithdraw(ctx context.Context, res *GetBalanceRes) {
	var (
		err          error
		errField     *errs.Error = nil
		symbol       string      = res.CurrencyConfig.Symbol
		pendingWDRaw string
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if pendingWDRaw, err = s.withdrawExchangeRepo.GetPendingWithdraw(symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPendingWithdraw)
		return
	}

	res.PendingWDCoin = util.RawToCoin(pendingWDRaw, 8)
	if res.PendingWDIdr, err = s.marketService.ConvertCoinToIdr(res.PendingWDCoin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPendingWithdraw)
	}
}
