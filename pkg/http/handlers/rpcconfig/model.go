package rpcconfig

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

type GetRes struct {
	RpcConfig domain.RpcConfig `json:"rpc_config"`
	Error     string           `json:"error"`
}

type ListRes struct {
	RpcConfigs []domain.RpcConfig `json:"rpc_configs"`
	Error      string             `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
