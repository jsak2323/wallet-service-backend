package btcxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type BtcService struct {
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.SystemConfigRepository
}

func (bs *BtcService) GetSymbol() string {
    return "BTC"
}

func (bs *BtcService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return bs.healthCheckRepo
}

func NewBtcService(healthCheckRepo hc.HealthCheckRepository, systemConfigRepo sc.SystemConfigRepository) *BtcService{
    return &BtcService{
        healthCheckRepo,
        systemConfigRepository,
    }
}