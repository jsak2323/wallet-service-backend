package currency

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	cr "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

const errInternalServer = "Internal server error"

type CurrencyConfigService struct {
	ccRepo cc.CurrencyConfigRepository
	crRepo cr.Repository
	rcRepo rc.RpcConfigRepository
}

func NewCurrencyConfigService(
		ccRepo cc.CurrencyConfigRepository,
		crRepo cr.Repository,
		rcRepo rc.RpcConfigRepository,
	) *CurrencyConfigService {
	return &CurrencyConfigService{
		ccRepo: ccRepo,
		crRepo: crRepo,
		rcRepo: rcRepo,
	}
}
