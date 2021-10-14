package rpcconfig

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

const errInternalServer = "Internal server error"

type RpcConfigService struct {
	rcRepo domain.RpcConfigRepository
}

func NewRpcConfigService(rcRepo domain.RpcConfigRepository) *RpcConfigService {
	return &RpcConfigService{rcRepo: rcRepo}
}
