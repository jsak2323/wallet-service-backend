package balance

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
)

type ListRes struct {
	Balances []WalletBalance `json:"balances"`
	Error  	 string		   `json:"error"`
}

type WalletBalance struct {
	cc.CurrencyConfig
	hw.WalletBalance
}
