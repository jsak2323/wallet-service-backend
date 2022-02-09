package rpcmethod

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	rcrm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type RpcMethodService interface {
	List(ctx context.Context, page, limit int) (resp []domain.RpcMethod, err error)
	GetByRpcConfigId(ctx context.Context, reqRpcConfigId int) (resp []domain.RpcMethod, err error)
	Create(ctx context.Context, req domain.RpcMethod) (err error)
	Update(ctx context.Context, req domain.UpdateRpcMethod) (err error)
	Delete(ctx context.Context, id, RpcConfigId int) (err error)
}

type rpcMethodService struct {
	validator util.CustomValidator
	rmRepo    rm.Repository
	rcrmRepo  rcrm.Repository
}

func NewRpcMethodService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *rpcMethodService {
	return &rpcMethodService{
		validator,
		mysqlRepos.RpcMethod,
		mysqlRepos.RpcConfigRpcMethod,
	}
}
