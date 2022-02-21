package wallet

import (
	"context"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *walletService) SetColdBalanceDetails(ctx context.Context, res *GetBalanceRes) {
	var (
		symbol   string           = res.CurrencyConfig.Symbol
		cbs      []cb.ColdBalance = s.coldWalletService.GetBalance(ctx, res.CurrencyConfig.Id)
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	for _, cb := range cbs {
		var coldBalanceDetail = BalanceDetail{Id: cb.Id, Name: cb.Name, Type: cb.Type}

		coldBalanceDetail.Coin = cb.Balance
		coldBalanceDetail.Address = cb.Address
		coldBalanceDetail.FireblocksName = cb.FireblocksName

		if coldBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(coldBalanceDetail.Coin, symbol); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		if res.TotalColdCoin, err = util.AddCurrency(res.TotalColdCoin, coldBalanceDetail.Coin); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		if res.TotalColdIdr, err = util.AddCurrency(res.TotalColdIdr, coldBalanceDetail.Idr); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		res.ColdBalances = append(res.ColdBalances, coldBalanceDetail)
	}
}
