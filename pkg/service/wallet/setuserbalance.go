package wallet

import (
	"context"

	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *walletService) SetUserBalanceDetails(ctx context.Context, res *GetBalanceRes) {
	var (
		tcb                 ub.TotalCoinBalance
		err                 error
		errField            *errs.Error   = nil
		symbol              string        = res.CurrencyConfig.Symbol
		frozenBalanceDetail BalanceDetail = BalanceDetail{Name: "Frozen"}
		liquidBalanceDetail BalanceDetail = BalanceDetail{Name: "Liquid"}
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if tcb, err = s.userBalanceRepo.GetTotalCoinBalance(symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	if liquidBalanceDetail.Coin = util.RawToCoin(tcb.Total, 8); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	} else if liquidBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(liquidBalanceDetail.Coin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	res.UserBalances = append(res.UserBalances, liquidBalanceDetail)

	frozenBalanceDetail.Coin = util.RawToCoin(tcb.TotalFrozen, 8)
	if frozenBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(frozenBalanceDetail.Coin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	res.UserBalances = append(res.UserBalances, frozenBalanceDetail)

	if res.TotalUserCoin, err = util.AddCurrency(liquidBalanceDetail.Coin, frozenBalanceDetail.Coin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	if res.TotalUserIdr, err = util.AddCurrency(liquidBalanceDetail.Idr, frozenBalanceDetail.Idr); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}
}
