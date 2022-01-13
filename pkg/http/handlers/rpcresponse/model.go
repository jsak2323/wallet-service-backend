package rpcresponse

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type GetRes struct {
	RpcResponse domain.RpcResponse `json:"rpc_response"`
	Error       *errs.Error        `json:"error"`
}

type ListRes struct {
	RpcResponses []domain.RpcResponse `json:"rpc_responses"`
	Error        *errs.Error          `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
