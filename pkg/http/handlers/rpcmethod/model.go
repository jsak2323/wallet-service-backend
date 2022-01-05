package rpcmethod

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type GetRes struct {
	RpcMethod domain.RpcMethod `json:"rpc_method"`
	Error     *errs.Error      `json:"error"`
}

type ListRes struct {
	RpcMethods []domain.RpcMethod `json:"rpc_methods"`
	Error      *errs.Error        `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
