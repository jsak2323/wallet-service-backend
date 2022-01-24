package rpcmethod

import (
	rcrmdomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	rmdomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type RpcMethodService struct {
	rmRepo    rmdomain.Repository
	rcrmRepo  rcrmdomain.Repository
	validator util.CustomValidator
}

func NewRpcMethodService(rmRepo rmdomain.Repository, rcrmRepo rcrmdomain.Repository, validator util.CustomValidator) *RpcMethodService {
	return &RpcMethodService{
		rmRepo:    rmRepo,
		rcrmRepo:  rcrmRepo,
		validator: validator,
	}
}
