package cron

import (
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

func runUpdateDeposit(sleep time.Duration, mysqlRepos mysql.MysqlRepositories) {
	fmt.Println("Initializing deposit service...")

	moduleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	listTransactionsService := h.NewListTransactionsService(moduleServices)
	depositService := hc.NewDepositService(moduleServices, listTransactionsService, mysqlRepos.Deposit)

	func() {
		for {
			fmt.Println()
			fmt.Println("- Running update deposit ...")
			depositService.Update()

			fmt.Println("- Finished running update deposit, sleeping for " + sleep.String() + " ...")
			countDownSleep("deposit", int(sleep.Seconds()))
			fmt.Println()
		}
	}()
}
