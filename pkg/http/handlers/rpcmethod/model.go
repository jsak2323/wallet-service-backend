package rpcmethod

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
)

type GetRes struct {
	RpcMethod domain.RpcMethod `json:"rpc_method"`
	Error     string           `json:"error"`
}

type ListRes struct {
	RpcMethods []domain.RpcMethod `json:"rpc_methods"`
	Error      string             `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
