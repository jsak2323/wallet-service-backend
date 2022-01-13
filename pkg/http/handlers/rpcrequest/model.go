package rpcrequest

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type GetRes struct {
	RpcRequest domain.RpcRequest `json:"rpc_request"`
	Error      *errs.Error       `json:"error"`
}

type ListRes struct {
	RpcRequests []domain.RpcRequest `json:"rpc_requests"`
	Error       *errs.Error         `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
