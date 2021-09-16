package currency

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
)

const errInternalServer = "Internal server error"

type CurrencyConfigService struct {
	ccRepo domain.CurrencyConfigRepository
}

func NewCurrencyConfigService(ccRepo domain.CurrencyConfigRepository) *CurrencyConfigService {
	return &CurrencyConfigService{ccRepo: ccRepo}
}
