package service

import (
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/service/coldwallet"
	"github.com/btcid/wallet-services-backend-go/pkg/service/cron"
	"github.com/btcid/wallet-services-backend-go/pkg/service/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/service/deposit"
	"github.com/btcid/wallet-services-backend-go/pkg/service/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/service/market"
	"github.com/btcid/wallet-services-backend-go/pkg/service/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/service/role"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/service/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/service/user"
	"github.com/btcid/wallet-services-backend-go/pkg/service/userwallet"
	"github.com/btcid/wallet-services-backend-go/pkg/service/wallet"
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
	ColdWallet  coldwallet.ColdWalletService
	Market      market.MarketService
	UserWallet  userwallet.UserWalletService
	Wallet      wallet.WalletService
	Cron        cron.CronService
}

func New(
	validator util.CustomValidator,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
) Service {
	moduleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	coldWalletServices := coldwallet.NewColdWalletService(validator, mysqlRepos)
	marketServices := market.NewMarketService(moduleServices, mysqlRepos, exchangeApiRepos)
	WalletServices := wallet.NewWalletService(moduleServices, mysqlRepos, exchangeApiRepos, coldWalletServices, marketServices)
	walletRpcServices := walletrpc.NewWalletRpcService(validator, moduleServices, mysqlRepos)

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
		UserWallet:  userwallet.NewUserWalet(validator, mysqlRepos),
		WalletRpc:   walletRpcServices,
		ColdWallet:  coldWalletServices,
		Market:      marketServices,
		Wallet:      WalletServices,
		Cron:        cron.NewCronService(moduleServices, mysqlRepos, exchangeApiRepos, coldWalletServices, marketServices, WalletServices, walletRpcServices),
	}
	return svc
}
