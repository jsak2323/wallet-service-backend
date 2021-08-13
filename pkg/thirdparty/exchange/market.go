package exchange

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

type exchangeMarketRepository struct {}

func NewExchangeMarketRepository() domain.Repository {
	return &exchangeMarketRepository{}
}
