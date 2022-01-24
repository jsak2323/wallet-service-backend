package rpcresponse

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type RpcResponseService struct {
	rrsRepo   domain.Repository
	validator util.CustomValidator
}

func NewRpcResponseService(rrsRepo domain.Repository, validator util.CustomValidator) *RpcResponseService {
	return &RpcResponseService{rrsRepo: rrsRepo, validator: validator}
}
