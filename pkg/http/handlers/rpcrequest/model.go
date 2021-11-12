package rpcrequest

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
)

type GetRes struct {
	RpcRequest domain.RpcRequest `json:"rpc_request"`
	Error      string            `json:"error"`
}

type ListRes struct {
	RpcRequests []domain.RpcRequest `json:"rpc_requests"`
	Error       string              `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
