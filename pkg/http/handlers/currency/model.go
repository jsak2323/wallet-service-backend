package currency

import (
    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
)

type ListRes struct {
    CurrencyConfigs  	[]cc.CurrencyConfig `json:"currency_configs"`
    Error           	string				`json:"error"`
}