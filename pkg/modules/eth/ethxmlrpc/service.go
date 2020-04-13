package ethxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type EthService struct {
    healthCheckRepo hc.HealthCheckRepository
}

func (es *EthService) GetSymbol() string {
    return "ETH"
}

func (es *EthService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return es.healthCheckRepo
}

func NewEthService(healthCheckRepo hc.HealthCheckRepository) *EthService{
    return &EthService{
        healthCheckRepo,
    }
}