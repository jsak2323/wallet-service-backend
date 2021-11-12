package rpcrequest

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
)

const errInternalServer = "Internal server error"

type RpcRequestService struct {
	rrqRepo domain.Repository
}

func NewRpcRequestService(rrqRepo domain.Repository) *RpcRequestService {
	return &RpcRequestService{rrqRepo: rrqRepo}
}
