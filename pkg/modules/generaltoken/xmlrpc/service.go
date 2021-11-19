package xmlrpc

import(
    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralTokenService struct {
    ParentSymbol     string
    Symbol           string
    TokenType        string
    healthCheckRepo  hc.Repository
    systemConfigRepo sc.Repository
}

func (gts *GeneralTokenService) GetSymbol() string {
    return gts.Symbol
}

func (gts *GeneralTokenService) GetParentSymbol() string {
    return gts.ParentSymbol
}

func (gts *GeneralTokenService) GetHealthCheckRepo() hc.Repository {
    return gts.healthCheckRepo
}

func NewGeneralTokenService(
    currencyConfig cc.CurrencyConfig,
    healthCheckRepo hc.Repository,
    systemConfigRepo sc.Repository,
) *GeneralTokenService {
    return &GeneralTokenService{
        currencyConfig.ParentSymbol,
        currencyConfig.Symbol,
        currencyConfig.TokenType,
        healthCheckRepo,
        systemConfigRepo,
    }
}


