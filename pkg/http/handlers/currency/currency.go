package currency

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	cr "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

const errInternalServer = "Internal server error"

type CurrencyConfigService struct {
	ccRepo cc.Repository
	crRepo cr.Repository
	rcRepo rc.Repository
}

func NewCurrencyConfigService(
		ccRepo cc.Repository,
		crRepo cr.Repository,
		rcRepo rc.Repository,
	) *CurrencyConfigService {
	return &CurrencyConfigService{
		ccRepo: ccRepo,
		crRepo: crRepo,
		rcRepo: rcRepo,
	}
}
