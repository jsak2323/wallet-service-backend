package exchange

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

func (r *exchangeMarketRepository) LastPriceBySymbol(symbol, trade string) (price string, err error) {
	sumRes, err := summaries()
	if err != nil {
		return "0", err
	}

	price24H, ok := sumRes.Prices24H[symbol+trade]
	if !ok {
		return "0", domain.ErrMarketTradeNotFound
	}
	
	return price24H, nil
}