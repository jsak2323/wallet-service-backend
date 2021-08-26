package wallet

import (
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
)

type GetBalanceRes struct {
	CurrencyConfig cc.CurrencyConfig `json:"currency_config"`

	ColdBalances []BalanceDetail `json:"cold_balances"`
	HotBalances  []BalanceDetail `json:"hot_balances"`
	UserBalances []BalanceDetail `json:"user_balances"`
	
	TotalColdCoin	string `json:"total_cold_coin"`
	TotalColdIdr 	string `json:"total_cold_idr"`
	TotalHotCoin 	string `json:"total_hot_coin"`
	TotalHotIdr 	string `json:"total_hot_idr"`
	TotalUserIdr 	string `json:"total_user_idr"`
	TotalUserCoin 	string `json:"total_user_coin"`
	PendingWDCoin	string `json:"pending_wd_coin"`
	PendingWDIdr	string `json:"pending_wd_idr"`

	HotLimits 		hl.HotLimit	`json:"hot_limits"`

	Error string `json:"error"`
}

type BalanceDetail struct {
	Id 	 	int 	`json:"id"`
	Name 	string 	`json:"name"`
	Type 	string 	`json:"type"`
	Address string 	`json:"address"`
	Coin 	string 	`json:"coin"`
	Idr  	string 	`json:"idr"`
}