package wallet

import (
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *WalletService) FormatWalletBalanceCurrency(walletBalance *GetBalanceRes) {
	for i := range walletBalance.ColdBalances {
		walletBalance.ColdBalances[i].Idr = util.FormatCurrency(walletBalance.ColdBalances[i].Idr)
		walletBalance.ColdBalances[i].Coin = util.FormatCurrency(walletBalance.ColdBalances[i].Coin)
	}

	walletBalance.TotalColdCoin = util.FormatCurrency(walletBalance.TotalColdCoin)
	walletBalance.TotalColdIdr = util.FormatCurrency(walletBalance.TotalColdIdr)

	for i := range walletBalance.HotBalances {
		walletBalance.HotBalances[i].Idr = util.FormatCurrency(walletBalance.HotBalances[i].Idr)
		walletBalance.HotBalances[i].Coin = util.FormatCurrency(walletBalance.HotBalances[i].Coin)
	}

	walletBalance.TotalHotCoin = util.FormatCurrency(walletBalance.TotalHotCoin)
	walletBalance.TotalHotIdr = util.FormatCurrency(walletBalance.TotalHotIdr)

	for i := range walletBalance.UserBalances {
		walletBalance.UserBalances[i].Idr = util.FormatCurrency(walletBalance.UserBalances[i].Idr)
		walletBalance.UserBalances[i].Coin = util.FormatCurrency(walletBalance.UserBalances[i].Coin)
	}

	walletBalance.TotalUserCoin = util.FormatCurrency(walletBalance.TotalUserCoin)
	walletBalance.TotalUserIdr = util.FormatCurrency(walletBalance.TotalUserIdr)
}