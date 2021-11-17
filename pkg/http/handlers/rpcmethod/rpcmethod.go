package rpcmethod

import (
	rcrmdomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	rmdomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
)

const errInternalServer = "Internal server error"

type RpcMethodService struct {
	rmRepo   rmdomain.Repository
	rcrmRepo rcrmdomain.Repository
}

func NewRpcMethodService(rmRepo rmdomain.Repository, rcrmRepo rcrmdomain.Repository) *RpcMethodService {
	return &RpcMethodService{
		rmRepo:   rmRepo,
		rcrmRepo: rcrmRepo,
	}
}
