package btc

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type BtcService struct {
    healthCheckRepo hc.HealthCheckRepository
}

func NewBtcService(healthCheckRepo hc.HealthCheckRepository) *BtcService{
    return &BtcService{
        healthCheckRepo,
    }
}