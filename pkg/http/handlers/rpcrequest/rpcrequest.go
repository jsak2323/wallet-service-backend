package rpcrequest

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type RpcRequestService struct {
	rrqRepo   domain.Repository
	validator util.CustomValidator
}

func NewRpcRequestService(rrqRepo domain.Repository, validator util.CustomValidator) *RpcRequestService {
	return &RpcRequestService{rrqRepo: rrqRepo, validator: validator}
}
