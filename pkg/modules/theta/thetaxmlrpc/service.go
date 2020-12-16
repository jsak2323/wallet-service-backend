package thetaxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
)

type ThetaService struct {
	Symbol 			string
    healthCheckRepo hc.HealthCheckRepository
}

func (ts *ThetaService) GetSymbol() string {
    return ts.Symbol
}

func (ts *ThetaService) GetParentSymbol() string {
    return "THETA"
}

func (ts *ThetaService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return ts.healthCheckRepo
}

func NewThetaService(symbol string, healthCheckRepo hc.HealthCheckRepository) *ThetaService{
    return &ThetaService{
    	symbol,
        healthCheckRepo,
    }
}


