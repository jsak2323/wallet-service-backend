package rpcconfig

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rcrmDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type RpcConfigService struct {
	rcRepo    domain.Repository
	rcrmRepo  rcrmDomain.Repository
	validator util.CustomValidator
}

func NewRpcConfigService(rcRepo domain.Repository, rcrmRepo rcrmDomain.Repository, validator util.CustomValidator) *RpcConfigService {
	return &RpcConfigService{
		rcRepo:    rcRepo,
		rcrmRepo:  rcrmRepo,
		validator: validator,
	}
}
