package xmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
)

type GeneralTokenService struct {
    ParentSymbol    string
    Symbol          string
    healthCheckRepo hc.HealthCheckRepository
}

func (gts *GeneralTokenService) GetSymbol() string {
    return gts.Symbol
}

func (gts *GeneralTokenService) GetParentSymbol() string {
    return gts.ParentSymbol
}

func (gts *GeneralTokenService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return gts.healthCheckRepo
}

func NewGeneralTokenService(parentSymbol string, symbol string, healthCheckRepo hc.HealthCheckRepository) *GeneralTokenService{
    return &GeneralTokenService{
        parentSymbol,
        symbol,
        healthCheckRepo,
    }
}


