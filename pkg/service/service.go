package service

import (
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/service/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/service/deposit"
	"github.com/btcid/wallet-services-backend-go/pkg/service/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/service/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/service/role"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/service/user"
	"github.com/btcid/wallet-services-backend-go/pkg/service/walletrpc"
	"github.com/btcid/wallet-services-backend-go/pkg/service/withdraw"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type Service struct {
	Permission  permission.PermissionService
	User        user.UserService
	Role        role.RoleService
	Currency    currency.CurrencyService
	RpcConfig   rpcconfig.RpcConfigService
	RpcMethod   rpcmethod.RpcMethodService
	RpcRequest  rpcrequest.RpcRequestService
	Deposit     deposit.DepositService
	Fireblocks  fireblocks.FireblocksService
	RpcResponse rpcresponse.RpcResponseService
	Withdraw    withdraw.WithdrawService
	WalletRpc   walletrpc.WalletRpcService
}

func New(
	validator util.CustomValidator,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
) Service {
	moduleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	svc := Service{
		Permission:  permission.NewPermissionService(validator, mysqlRepos),
		User:        user.NewUserService(validator, mysqlRepos),
		Role:        role.NewRoleService(validator, mysqlRepos),
		Currency:    currency.NewCurrencyService(validator, mysqlRepos),
		RpcConfig:   rpcconfig.NewRpcConfigService(validator, mysqlRepos),
		RpcMethod:   rpcmethod.NewRpcMethodService(validator, mysqlRepos),
		RpcRequest:  rpcrequest.NewRpcRequestService(validator, mysqlRepos),
		Deposit:     deposit.NewDepositService(validator, mysqlRepos),
		Fireblocks:  fireblocks.NewFireblocksService(validator, mysqlRepos),
		RpcResponse: rpcresponse.NewRpcResponseService(validator, mysqlRepos),
		Withdraw:    withdraw.NewWithdrawService(validator, mysqlRepos),
		WalletRpc:   walletrpc.NewWalletRpcService(validator, moduleServices, mysqlRepos),
	}
	return svc
}
