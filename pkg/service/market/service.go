package market

import (
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type MarketService interface {
	ConvertCoinToIdr(amount string, symbol string) (string, error)
	ConvertIdrToCoin(amount string, symbol string) (string, error)
}

type marketService struct {
	marketRepo market.Repository
}

func NewMarketService(
	moduleServices *modules.ModuleServiceMap,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
) *marketService {
	return &marketService{
		exchangeApiRepos.Market,
	}
}
