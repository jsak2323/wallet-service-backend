package rpcrequest

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type RpcRequestService interface {
	GetByRpcMethodId(ctx context.Context, reqRpcMethodId int) (resp []domain.RpcRequest, err error)
	Create(ctx context.Context, req domain.RpcRequest) (err error)
	Update(ctx context.Context, req domain.UpdateRpcRequest) (err error)
	Delete(ctx context.Context, id, RpcMethodId int) (err error)
}

type rpcRequestService struct {
	validator util.CustomValidator
	rrqRepo   rr.Repository
}

func NewRpcRequestService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *rpcRequestService {
	return &rpcRequestService{
		validator,
		mysqlRepos.RpcRequest,
	}
}
