package exchange

import (
	"errors"
)

func (r *exchangeMarketRepository) LastPriceBySymbol(symbol, trade string) (price string, err error) {
	sumRes, err := summaries()
	if err != nil {
		return "0", err
	}

	ticker, ok := sumRes.Tickers[symbol+"_"+trade]
	if !ok {
		return "0", errors.New("market trade not found")
	}
	
	return ticker.Last, nil
}