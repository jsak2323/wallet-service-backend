package wallet

import (
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *WalletService) FormatWalletBalanceCurrency(walletBalance GetBalanceRes) (result GetBalanceRes) {

	result.CurrencyConfig = walletBalance.CurrencyConfig
	result.ColdBalances = make([]BalanceDetail, len(walletBalance.ColdBalances))
	result.HotBalances = make([]BalanceDetail, len(walletBalance.HotBalances))
	result.UserBalances = make([]BalanceDetail, len(walletBalance.UserBalances))

	for i := range walletBalance.ColdBalances {
		result.ColdBalances[i] = walletBalance.ColdBalances[i]
		result.ColdBalances[i].Idr = util.FormatCurrency(walletBalance.ColdBalances[i].Idr)
		result.ColdBalances[i].Coin = util.FormatCurrency(walletBalance.ColdBalances[i].Coin)
	}

	result.TotalColdCoin = util.FormatCurrency(walletBalance.TotalColdCoin)
	result.TotalColdIdr = util.FormatCurrency(walletBalance.TotalColdIdr)

	for i := range walletBalance.HotBalances {
		result.HotBalances[i] = walletBalance.HotBalances[i]
		result.HotBalances[i].Idr = util.FormatCurrency(walletBalance.HotBalances[i].Idr)
		result.HotBalances[i].Coin = util.FormatCurrency(walletBalance.HotBalances[i].Coin)
	}

	result.TotalHotCoin = util.FormatCurrency(walletBalance.TotalHotCoin)
	result.TotalHotIdr = util.FormatCurrency(walletBalance.TotalHotIdr)
	
	for i := range walletBalance.UserBalances {
		result.UserBalances[i] = walletBalance.UserBalances[i]
		result.UserBalances[i].Idr = util.FormatCurrency(walletBalance.UserBalances[i].Idr)
		result.UserBalances[i].Coin = util.FormatCurrency(walletBalance.UserBalances[i].Coin)
	}

	result.TotalUserCoin = util.FormatCurrency(walletBalance.TotalUserCoin)
	result.TotalUserIdr = util.FormatCurrency(walletBalance.TotalUserIdr)
	result.PendingWDCoin = util.FormatCurrency(walletBalance.PendingWDCoin)
	result.PendingWDIdr = util.FormatCurrency(walletBalance.PendingWDIdr)
	result.HotPercent = walletBalance.HotPercent
	result.HotColdPercent = walletBalance.HotColdPercent

	return result
}