package user

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
)

type TotalUserBalanceRes struct {
	TokenTypes []cc.CurrencyConfig
	Balance ub.TotalCoinBalance
}