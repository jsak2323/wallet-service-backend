package market

import (
	"errors"
	"math/big"
	"strings"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *marketService) ConvertCoinToIdr(amount string, symbol string) (string, error) {
	lastPrice, err := s.getLastPrice(symbol)
	if err != nil {
		return "0", errs.AddTrace(err)
	}

	if lastPrice == "" {
		lastPrice = "0"
	}
	if amount == "" {
		amount = "0"
	}

	lastPriceFloat, ok := big.NewFloat(0).SetString(lastPrice)
	if !ok {
		return "", errs.AddTrace(errors.New("Fail parsing lastPrice"))
	}

	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return "", errs.AddTrace(errors.New("Fail parsing amount"))
	}

	return amountFloat.Mul(lastPriceFloat, amountFloat).Text('f', 0), nil
}

func (s *marketService) ConvertIdrToCoin(amount string, symbol string) (string, error) {
	lastPrice, err := s.getLastPrice(symbol)
	if err != nil {
		return "0", err
	}

	if lastPrice == "" {
		lastPrice = "0"
	}
	if amount == "" {
		amount = "0"
	}

	lastPriceFloat, ok := big.NewFloat(0).SetString(lastPrice)
	if !ok {
		return "", errs.AddTrace(errors.New("Fail parsing lastPrice"))
	}

	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return "", errs.AddTrace(errors.New("Fail parsing amount"))
	}

	return amountFloat.Quo(amountFloat, lastPriceFloat).Text('f', 8), nil
}

func (s *marketService) getLastPrice(symbol string) (price string, err error) {
	// TODO caching

	if price, err = s.marketRepo.LastPriceBySymbol(strings.ToLower(symbol), "idr"); err != nil {

		// check symbolusdt table,
		// assuming that no symbolidr table actually means symbol is not in idr market
		// TODO make sure the assumption is right
		if err == market.ErrMarketTradeNotFound {
			if price, err = s.marketRepo.LastPriceBySymbol(strings.ToLower(symbol), "usdt"); err != nil {
				return "0", errs.AddTrace(err)
			}
		}

		return "0", err
	}

	return price, nil
}
