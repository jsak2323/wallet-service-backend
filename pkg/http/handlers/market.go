package handlers

import (
	"math/big"
	"strings"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

type MarketService struct {
	marketRepo market.Repository
	LastPrices map[string]string
}

func NewMarketService(marketRepo market.Repository) *MarketService {
	return &MarketService{
		marketRepo: marketRepo,
	}
}

func (s *MarketService) ConvertCoinToIdr(amount string, symbol string) (string, error) {
	lastPrice, err := s.getLastPrice(symbol)
	if err != nil {
		return "0", err
	}
	
	lastPriceFloat, _ := big.NewFloat(0).SetString(lastPrice)
	
	amountFloat, _ := big.NewFloat(0).SetString(amount)
	
	return amountFloat.Mul(lastPriceFloat, amountFloat).Text('f', 0), nil
}

func (s *MarketService) ConvertIdrToCoin(amount string, symbol string) (string, error) {
	lastPrice, err := s.getLastPrice(symbol)
	if err != nil {
		return "0", err
	}

	lastPriceFloat, _ := big.NewFloat(0).SetString(lastPrice)

	amountFloat, _ := big.NewFloat(0).SetString(amount)

	return amountFloat.Quo(amountFloat, lastPriceFloat).Text('f', 8), nil
}

func (s *MarketService) getLastPrice(symbol string) (price string, err error) {
	// TODO caching
	
	if price, err = s.marketRepo.LastPriceBySymbol(strings.ToLower(symbol+"idr")); err != nil {

		// check symbolusdt table,
		// assuming that no symbolidr table actually means symbol is not in idr market
		// TODO make sure the assumption is right
		if strings.Contains(err.Error(), "1146") {
			if price, err = s.marketRepo.LastPriceBySymbol(strings.ToLower(symbol+"usdt")); err != nil {
				return "0", err
			}
		}

		return "0", err
	}
	
	return price, nil
}
