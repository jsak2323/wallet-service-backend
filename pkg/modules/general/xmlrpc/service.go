package xmlrpc

import(
    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralService struct {
    Symbol           string
    CurrencyConfigId int
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.SystemConfigRepository
}

func (gs *GeneralService) GetSymbol() string {
    return gs.Symbol
}

func (gs *GeneralService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return gs.healthCheckRepo
}

func NewGeneralService(
    currencyConfig cc.CurrencyConfig,
    healthCheckRepo hc.HealthCheckRepository,
    systemConfigRepo sc.SystemConfigRepository,
) *GeneralService {
    return &GeneralService{
        currencyConfig.Symbol,
        currencyConfig.Id,
        healthCheckRepo,
        systemConfigRepo,
    }
}


