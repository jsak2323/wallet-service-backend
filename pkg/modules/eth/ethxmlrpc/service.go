package ethxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type EthService struct {
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.SystemConfigRepository
}

func (es *EthService) GetSymbol() string {
    return "ETH"
}

func (es *EthService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return es.healthCheckRepo
}

func NewEthService(healthCheckRepo hc.HealthCheckRepository, systemConfigRepo sc.SystemConfigRepository) *EthService{
    return &EthService{
        healthCheckRepo,
        systemConfigRepo,
    }
}