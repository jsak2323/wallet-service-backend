package rpcresponse

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type GetRes struct {
	RpcResponse domain.RpcResponse `json:"rpc_response"`
	Error       string             `json:"error"`
}

type ListRes struct {
	RpcResponses []domain.RpcResponse `json:"rpc_responses"`
	Error        string               `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
