package ethxmlrpc

import (
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type EthService struct {
    currencyConfigId int
    healthCheckRepo  hc.Repository
    systemConfigRepo sc.Repository
}

func (es *EthService) GetSymbol() string {
    return "ETH"
}

func (es *EthService) GetHealthCheckRepo() hc.Repository {
    return es.healthCheckRepo
}

func NewEthService(currencyConfigId int, healthCheckRepo hc.Repository, systemConfigRepo sc.Repository) *EthService{
    return &EthService{
        currencyConfigId,
        healthCheckRepo,
        systemConfigRepo,
    }
}