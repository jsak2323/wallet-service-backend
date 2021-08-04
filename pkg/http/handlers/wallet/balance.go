package wallet

import (
	"strings"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	modulesm "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type WalletBalance struct {
	ColdBalances []BalanceDetail `json:"cold_balances"`
	HotBalances  []BalanceDetail `json:"hot_balances"`
	UserBalances []BalanceDetail `json:"user_balances"`
	
	TotalColdCoin, TotalNodeCoin string
	TotalUserCoin string `json:"total_user_coin"`
	TotalColdIdr, TotalNodeIdr string
	TotalUserIdr string `json:"total_user_idr"`
}

type BalanceDetail struct {
	Id 	 int 	 `json:"id"`
	Name string	 `json:"name"`
	Coin string  `json:"coin"`
	Idr  string	 `json:"idr"`
}

func (s *WalletService) SetColdBalanceDetails(symbol string, walletBalance *WalletBalance) {
	var (
		cbs    []cb.ColdBalance = s.coldWalletService.GetBalance(symbol)
		err    error
	)
	
	for _, cb := range cbs {
		var coldBalanceDetail = BalanceDetail{ Id: cb.Id, Name: cb.Type }

		coldBalanceDetail.Coin = cb.Balance
		
		if coldBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(coldBalanceDetail.Coin, symbol); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails ConvertCoinToIdr("+cb.Name+", "+cb.Balance+") err: "+err.Error())
		}

		if walletBalance.TotalColdCoin, err = util.AddCoin(walletBalance.TotalColdCoin, coldBalanceDetail.Coin); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails AddCoin("+cb.Name+", "+cb.Balance+") err: "+err.Error())
		}

		if walletBalance.TotalColdIdr, err = util.AddIdr(walletBalance.TotalColdIdr, coldBalanceDetail.Idr); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails AddIdr("+cb.Name+", "+cb.Balance+") err: "+err.Error())
		}

		walletBalance.ColdBalances = append(walletBalance.ColdBalances, coldBalanceDetail)
	}
}

func (s *WalletService) SetHotBalanceDetails(symbol string, rpcConfigs []rc.RpcConfig, walletBalance *WalletBalance) {
	var err error
	
	for _, rpcConfig := range rpcConfigs {
		var hotBalanceDetail BalanceDetail = BalanceDetail{ Name: rpcConfig.Type }
		var res 			 *modulesm.GetBalanceRpcRes

		if res, err = (*s.moduleServices)[symbol].GetBalance(rpcConfig); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails node.GetBalance("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		hotBalanceDetail.Coin = res.Balance

		if hotBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(hotBalanceDetail.Coin, symbol); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails node.ConvertCoinToIdr("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		if walletBalance.TotalNodeCoin, err = util.AddCoin(walletBalance.TotalNodeCoin, hotBalanceDetail.Coin); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails AddCoin("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		if walletBalance.TotalNodeIdr, err = util.AddIdr(walletBalance.TotalNodeIdr, hotBalanceDetail.Idr); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails AddIdr("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		walletBalance.HotBalances = append(walletBalance.HotBalances, hotBalanceDetail)
	}
}

func (s *WalletService) SetUserBalanceDetails(symbol string, walletBalance *WalletBalance) {
	var tcb    ub.TotalCoinBalance
	var err    error
	var frozenBalanceDetail BalanceDetail = BalanceDetail{ Name: "frozen" }
	var liquidBalanceDetail BalanceDetail = BalanceDetail{ Name: "liquid" }
	
	if tcb, err = s.userBalanceRepo.GetTotalCoinBalance(symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails ub.GetTotalCoinBalance("+symbol+") err: "+err.Error())
	}

	if liquidBalanceDetail.Coin, err = util.RawToCoin(tcb.Total, 8); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails RawToCoin("+symbol+") err: "+err.Error())
	}

	if liquidBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(liquidBalanceDetail.Coin, symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails liquid.ConvertCoinToIdr("+symbol+") err: "+err.Error())
	}
	
	walletBalance.UserBalances = append(walletBalance.UserBalances, liquidBalanceDetail)
	
	if frozenBalanceDetail.Coin, err = util.RawToCoin(tcb.TotalFrozen, 8); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails RawToCoin("+symbol+") err: "+err.Error())
	}

	if frozenBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(frozenBalanceDetail.Coin, symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails frozen.ConvertCoinToIdr("+symbol+") err: "+err.Error())
	}

	walletBalance.UserBalances = append(walletBalance.UserBalances, frozenBalanceDetail)

	if walletBalance.TotalUserCoin, err = util.AddCoin(liquidBalanceDetail.Coin, frozenBalanceDetail.Coin); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails AddCoin("+symbol+") err: "+err.Error())
	}

	if walletBalance.TotalUserIdr, err = util.AddIdr(liquidBalanceDetail.Idr, frozenBalanceDetail.Idr); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails AddIdr("+symbol+") err: "+err.Error())
	}
}

func (s *WalletService) FormatWalletBalanceCurrency(symbol string, walletBalance *WalletBalance) {
	symbol = strings.ToUpper(symbol)

	for i := range walletBalance.ColdBalances {
		walletBalance.ColdBalances[i].Idr = util.FormatCurrency(walletBalance.ColdBalances[i].Idr, "IDR")
		// walletBalance.ColdBalances[i].Coin = util.FormatCurrency(walletBalance.ColdBalances[i].Coin, symbol)
	}

	// walletBalance.TotalColdCoin = util.FormatCurrency(walletBalance.TotalColdCoin, symbol)
	walletBalance.TotalColdIdr = util.FormatCurrency(walletBalance.TotalColdIdr, "IDR")

	for i := range walletBalance.HotBalances {
		walletBalance.HotBalances[i].Idr = util.FormatCurrency(walletBalance.HotBalances[i].Idr, "IDR")
		// walletBalance.HotBalances[i].Coin = util.FormatCurrency(walletBalance.HotBalances[i].Coin, symbol)
	}

	// walletBalance.TotalNodeCoin = util.FormatCurrency(walletBalance.TotalNodeCoin, symbol)
	walletBalance.TotalNodeIdr = util.FormatCurrency(walletBalance.TotalNodeIdr, "IDR")

	for i := range walletBalance.UserBalances {
		walletBalance.UserBalances[i].Idr = util.FormatCurrency(walletBalance.UserBalances[i].Idr, "IDR")
		// walletBalance.UserBalances[i].Coin = util.FormatCurrency(walletBalance.UserBalances[i].Coin, symbol)
	}

	// walletBalance.TotalUserCoin = util.FormatCurrency(walletBalance.TotalUserCoin, symbol)
	walletBalance.TotalUserIdr = util.FormatCurrency(walletBalance.TotalUserIdr, "IDR")
}