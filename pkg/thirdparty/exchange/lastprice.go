package exchange

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/market"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (r *exchangeMarketRepository) LastPriceBySymbol(symbol, trade string) (price string, err error) {
	sumRes, err := summaries()
	if err != nil {
		return "0", errs.AddTrace(err)
	}

	price24H, ok := sumRes.Prices24H[symbol+trade]
	if !ok {
		return "0", errs.AddTrace(domain.ErrMarketTradeNotFound)
	}

	return price24H, nil
}
