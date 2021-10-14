package wallet

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
)

type GetBalanceRes struct {
	CurrencyConfig cc.CurrencyConfig `json:"currency_config"`

	ColdBalances []BalanceDetail `json:"cold_balances"`
	HotBalances  []BalanceDetail `json:"hot_balances"`
	UserBalances []BalanceDetail `json:"user_balances"`

	TotalColdCoin  string `json:"total_cold_coin"`
	TotalColdIdr   string `json:"total_cold_idr"`
	TotalHotCoin   string `json:"total_hot_coin"`
	TotalHotIdr    string `json:"total_hot_idr"`
	TotalUserIdr   string `json:"total_user_idr"`
	TotalUserCoin  string `json:"total_user_coin"`
	PendingWDCoin  string `json:"pending_wd_coin"`
	PendingWDIdr   string `json:"pending_wd_idr"`
	HotPercent     string `json:"hot_percent"`
	HotColdPercent string `json:"hot_cold_percent"`

	HotLimits hl.HotLimit `json:"hot_limits"`

	Error string `json:"error"`
}

type BalanceDetail struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Address        string `json:"address"`
	FireblocksName string `json:"fireblocks_name"`
	Coin           string `json:"coin"`
	Idr            string `json:"idr"`
}
