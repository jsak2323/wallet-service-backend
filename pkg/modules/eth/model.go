package eth

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type EthService struct {
    healthCheckRepo hc.HealthCheckRepository
}

func NewEthService(healthCheckRepo hc.HealthCheckRepository) *EthService{
    return &EthService{
        healthCheckRepo,
    }
}