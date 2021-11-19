package rpcconfig

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rcrmDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
)

const errInternalServer = "Internal server error"

type RpcConfigService struct {
	rcRepo   domain.Repository
	rcrmRepo rcrmDomain.Repository
}

func NewRpcConfigService(rcRepo domain.Repository, rcrmRepo rcrmDomain.Repository) *RpcConfigService {
	return &RpcConfigService{
		rcRepo: rcRepo,
		rcrmRepo: rcrmRepo,
	}
}
