package currency

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
)

type ListRes struct {
	CurrencyConfigs []domain.CurrencyConfig `json:"currency_configs"`
	Error           string                  `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
