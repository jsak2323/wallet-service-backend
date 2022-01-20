package cron

import (
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func runUpdateWithdraw(sleep time.Duration, mysqlRepos mysql.MysqlRepositories, exchangeApiRepos exchange.APIRepositories) {
	fmt.Println("Initializing withdraw service...")

	moduleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	marketService := h.NewMarketService(exchangeApiRepos.Market)
	listWithdrawsService := h.NewListWithdrawsService(moduleServices)
	withdrawService := hc.NewWithdrawService(moduleServices, listWithdrawsService, marketService, mysqlRepos.Withdraw)

	func() {
		for {
			fmt.Println()
			fmt.Println("- Running update withdraw ...")
			withdrawService.Update()

			fmt.Println("- Finished running update withdraw, sleeping for " + sleep.String() + " ...")
			countDownSleep("withdraw", int(sleep.Seconds()))
			fmt.Println()
		}
	}()
}
