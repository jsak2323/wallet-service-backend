package btcxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type BtcService struct {
    healthCheckRepo hc.HealthCheckRepository
}

func (bs *BtcService) GetSymbol() string {
    return "BTC"
}

func (bs *BtcService) GetHealthCheckRepo() hc.HealthCheckRepository {
    return bs.healthCheckRepo
}

func NewBtcService(healthCheckRepo hc.HealthCheckRepository) *BtcService{
    return &BtcService{
        healthCheckRepo,
    }
}