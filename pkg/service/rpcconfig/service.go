package rpcconfig

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rcrm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type RpcConfigService interface {
	List(ctx context.Context, page, limit int) (resp []domain.RpcConfig, err error)
	GetById(ctx context.Context, idRpcConfig int) (res domain.RpcConfig, err error)
	Activate(ctx context.Context, id int) (err error)
	Deactivate(ctx context.Context, idRpcConfig int) (err error)
	Create(ctx context.Context, req domain.RpcConfig) (err error)
	CreateRpcMethod(ctx context.Context, req handlerRpcConfig.RpcConfigRpcMethodReq) (err error)
	Update(ctx context.Context, req domain.UpdateRpcConfig) (err error)
	DeleteRpcMethod(ctx context.Context, roleId, permissionId int) (err error)
}

type rpcConfigService struct {
	validator util.CustomValidator
	rcRepo    rc.Repository
	rcrmRepo  rcrm.Repository
}

func NewRpcConfigService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *rpcConfigService {
	return &rpcConfigService{
		validator,
		mysqlRepos.RpcConfig,
		mysqlRepos.RpcConfigRpcMethod,
	}
}
