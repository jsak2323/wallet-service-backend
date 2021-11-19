package cron

import (
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func runCheckBalance(sleep time.Duration, mysqlRepos mysql.MysqlRepositories, exchangeApiRepos exchange.APIRepositories) {
	fmt.Println("Initializing checkbalance service...")

	ModuleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	MarketService := h.NewMarketService(exchangeApiRepos.Market)
	coldWalletService := hcw.NewColdWalletService(mysqlRepos.ColdBalance)
	walletService := hw.NewWalletService(ModuleServices, coldWalletService, MarketService, mysqlRepos.Withdraw, exchangeApiRepos.HotLimit, mysqlRepos.UserBalance)
	checkBalanceService := hc.NewCheckBalanceService(walletService, coldWalletService, MarketService, ModuleServices, exchangeApiRepos.HotLimit, mysqlRepos.User)

	func() {
		for {
			fmt.Println("- Running checkbalance ...")
			checkBalanceService.CheckBalanceHandler()

			fmt.Println("- Finished running checkbalance, sleeping for " + sleep.String() + " ...")
			countDownSleep("checkbalance", int(sleep.Seconds()))
		}
	}()
}
