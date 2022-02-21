package cron

import (
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

func runHealthCheck(sleep time.Duration, mysqlRepos mysql.MysqlRepositories) {
	fmt.Println("Initializing healthcheck service...")

	moduleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	healthCheckService := hc.NewHealthCheckService(moduleServices, mysqlRepos.HealthCheck, mysqlRepos.SystemConfig)

	func() {
		for {
			fmt.Println()
			fmt.Println("- Running healthcheck ...")
			healthCheckService.HealthCheckHandler()

			fmt.Println("- Finished running healthcheck, sleeping for " + sleep.String() + " ...")
			countDownSleep("healthcheck", int(sleep.Seconds()))
			fmt.Println()
		}
	}()
}
