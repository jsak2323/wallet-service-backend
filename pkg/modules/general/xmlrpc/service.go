package xmlrpc

import(
    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralService struct {
    Symbol           string
    CurrencyConfigId int
    healthCheckRepo  hc.Repository
    systemConfigRepo sc.Repository
}

func (gs *GeneralService) GetSymbol() string {
    return gs.Symbol
}

func (gs *GeneralService) GetHealthCheckRepo() hc.Repository {
    return gs.healthCheckRepo
}

func NewGeneralService(
    currencyConfig cc.CurrencyConfig,
    healthCheckRepo hc.Repository,
    systemConfigRepo sc.Repository,
) *GeneralService {
    return &GeneralService{
        currencyConfig.Symbol,
        currencyConfig.Id,
        healthCheckRepo,
        systemConfigRepo,
    }
}


