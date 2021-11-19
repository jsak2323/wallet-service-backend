package exchange

import (
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

type APIRepositories struct {
	HotLimit hl.Repository
	Market   market.Repository
}

func NewAPIRepositories() APIRepositories {
	return APIRepositories{
		HotLimit: NewExchangeHotLimitRepository(),
		Market:   NewExchangeMarketRepository(),
	}
}
