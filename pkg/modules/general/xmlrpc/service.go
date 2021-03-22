package xmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralService struct {
    Symbol           string
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.systemConfigRepository
}

func (gs *GeneralService) GetSymbol() string {
    return gs.Symbol
}

func (gs *GeneralService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return gs.healthCheckRepo
}

func NewGeneralService(
    symbol string, 
    healthCheckRepo hc.HealthCheckRepository,
    systemConfigRepo sc.SystemConfigRepository
) *GeneralService {
    return &GeneralService{
        symbol,
        healthCheckRepo,
        systemConfigRepo,
    }
}


