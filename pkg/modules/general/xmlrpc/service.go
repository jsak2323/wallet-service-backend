package xmlrpc

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type GeneralService struct {
    Symbol          string
    healthCheckRepo hc.HealthCheckRepository
}

func (gs *GeneralService) GetSymbol() string {
    return gs.Symbol
}

func (gs *GeneralService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return gs.healthCheckRepo
}

func NewGeneralService(symbol string, healthCheckRepo hc.HealthCheckRepository) *GeneralService {
    return &GeneralService{
        symbol,
        healthCheckRepo,
    }
}