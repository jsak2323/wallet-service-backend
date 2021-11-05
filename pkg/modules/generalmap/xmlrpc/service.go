package xmlrpc

import (
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

type GeneralMapService struct {
	ParentSymbol     string
	Symbol           string
	healthCheckRepo  hc.HealthCheckRepository
	systemConfigRepo sc.SystemConfigRepository
	rpcMethodRepo    rm.Repository
	rpcRequestRepo   rrq.Repository
	rpcResponseRepo  rrs.Repository
}

func (gms *GeneralMapService) GetSymbol() string {
	return gms.Symbol
}

func (gms *GeneralMapService) GetParentSymbol() string {
	return gms.ParentSymbol
}

func (gms *GeneralMapService) GetHealthCheckRepo() hc.HealthCheckRepository {
	return gms.healthCheckRepo
}

func NewGeneralMapService(
	parentSymbol string,
	symbol string,
	healthCheckRepo hc.HealthCheckRepository,
	systemConfigRepo sc.SystemConfigRepository,
	rpcMethodRepo rm.Repository,
	rpcRequestRepo rrq.Repository,
	rpcResponsRepo rrs.Repository,
) *GeneralMapService {
	return &GeneralMapService{
		parentSymbol,
		symbol,
		healthCheckRepo,
		systemConfigRepo,
		rpcMethodRepo,
		rpcRequestRepo,
		rpcResponsRepo,
	}
}
