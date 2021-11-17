package rpcresponse

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

const errInternalServer = "Internal server error"

type RpcResponseService struct {
	rrsRepo domain.Repository
}

func NewRpcResponseService(rrsRepo domain.Repository) *RpcResponseService {
	return &RpcResponseService{rrsRepo: rrsRepo}
}
