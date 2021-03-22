package xmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralTokenService struct {
    ParentSymbol     string
    Symbol           string
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.SystemConfigRepository
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

func NewGeneralTokenService(
    parentSymbol string, 
    symbol string, 
    healthCheckRepo hc.HealthCheckRepository,
    systemConfigRepo sc.SystemConfigRepository,
) *GeneralTokenService {
    return &GeneralTokenService{
        parentSymbol,
        symbol,
        healthCheckRepo,
        systemConfigRepo,
    }
}


