package mysql

import (
	"database/sql"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
)

type MysqlRepositories struct {
	ColdBalance        coldbalance.Repository
	CurrencyConfig     currencyconfig.CurrencyConfigRepository
	CurrencyRpc        currencyrpc.Repository
	HealthCheck        healthcheck.HealthCheckRepository
	Permission         permission.Repository
	Role               role.Repository
	RolePermission     rolepermission.Repository
	RpcConfig          rpcconfig.RpcConfigRepository
	RpcConfigRpcMethod rpcconfigrpcmethod.Repository
	RpcMethod          rpcmethod.Repository
	RpcRequest         rpcrequest.Repository
	RpcResponse        rpcresponse.Repository
	SystemConfig       systemconfig.SystemConfigRepository
	User               user.Repository
	UserBalance        userbalance.Repository
	UserRole           userrole.Repository
	Withdraw           withdraw.Repository
}

func NewMysqlRepositories(localDB *sql.DB, exchangeSlaveDB *sql.DB) MysqlRepositories {
	return MysqlRepositories{
		ColdBalance:        NewMysqlColdBalanceRepository(localDB),
		CurrencyConfig:     NewMysqlCurrencyConfigRepository(localDB),
		CurrencyRpc:        NewMysqlCurrencyRpcRepository(localDB),
		HealthCheck:        NewMysqlHealthCheckRepository(localDB),
		Permission:         NewMysqlPermissionRepository(localDB),
		Role:               NewMysqlRoleRepository(localDB),
		RolePermission:     NewMysqlRolePermissionRepository(localDB),
		RpcConfig:          NewMysqlRpcConfigRepository(localDB),
		RpcConfigRpcMethod: NewMysqlRpcConfigRpcMethodRepository(localDB),
		RpcMethod:          NewMysqlRpcMethodRepository(localDB),
		RpcRequest:         NewMysqlRpcRequestRepository(localDB),
		RpcResponse:        NewMysqlRpcResponseRepository(localDB),
		SystemConfig:       NewMysqlSystemConfigRepository(localDB),
		User:               NewMysqlUserRepository(localDB),
		UserBalance:        NewMysqlUserBalanceRepository(exchangeSlaveDB),
		UserRole:           NewMysqlUserRoleRepository(localDB),
		Withdraw:           NewMysqlWithdrawRepository(exchangeSlaveDB),
	}
}
