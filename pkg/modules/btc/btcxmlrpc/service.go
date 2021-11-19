package btcxmlrpc

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type BtcService struct {
    currencyConfigId int
    healthCheckRepo  hc.Repository
    systemConfigRepo sc.Repository
}

func (bs *BtcService) GetSymbol() string {
    return "BTC"
}

func (bs *BtcService) GetHealthCheckRepo() hc.Repository {
    return bs.healthCheckRepo
}

func NewBtcService(currencyConfigId int, healthCheckRepo hc.Repository, systemConfigRepo sc.Repository) *BtcService{
    return &BtcService{
        currencyConfigId,
        healthCheckRepo,
        systemConfigRepo,
    }
}