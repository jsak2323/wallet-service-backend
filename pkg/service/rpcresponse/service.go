package rpcresponse

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type RpcResponseService interface {
	GetByRpcMethod(ctx context.Context, reqRpcMethodId int) (resp []rpcresponse.RpcResponse, err error)
	Create(ctx context.Context, req domain.CreateRpcResponse) (err error)
	Update(ctx context.Context, req domain.RpcResponse) (err error)
	Delete(ctx context.Context, id, rpcMethodId int) (err error)
}

type rpcResponseService struct {
	validator       util.CustomValidator
	rpcresponseRepo rpcresponse.Repository
}

func NewRpcResponseService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *rpcResponseService {
	return &rpcResponseService{
		validator,
		mysqlRepos.RpcResponse,
	}
}
