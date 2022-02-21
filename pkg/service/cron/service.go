package cron

import (
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	d "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	wd "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	cws "github.com/btcid/wallet-services-backend-go/pkg/service/coldwallet"
	ms "github.com/btcid/wallet-services-backend-go/pkg/service/market"
	ws "github.com/btcid/wallet-services-backend-go/pkg/service/wallet"
	wrs "github.com/btcid/wallet-services-backend-go/pkg/service/walletrpc"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type CronService interface {
	CheckBalanceHandler()
	UpdateDeposit()
	HealthCheckHandler()
	UpdateWithdraw()
}

type cronService struct {
	moduleServices    *modules.ModuleServiceMap
	depositRepo       d.Repository
	withdrawRepo      wd.Repository
	hotLimitRepo      hl.Repository
	userRepo          user.Repository
	healthCheckRepo   hc.Repository
	systemConfigRepo  sc.Repository
	coldWalletService cws.ColdWalletService
	marketService     ms.MarketService
	walletService     ws.WalletService
	walletRpcService  wrs.WalletRpcService
}

func NewCronService(
	moduleServices *modules.ModuleServiceMap,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
	coldWalletService cws.ColdWalletService,
	marketService ms.MarketService,
	walletService ws.WalletService,
	walletRpcService wrs.WalletRpcService,
) *cronService {
	return &cronService{
		moduleServices,
		mysqlRepos.Deposit,
		mysqlRepos.Withdraw,
		exchangeApiRepos.HotLimit,
		mysqlRepos.User,
		mysqlRepos.HealthCheck,
		mysqlRepos.SystemConfig,
		coldWalletService,
		marketService,
		walletService,
		walletRpcService,
	}
}
