package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func main() {
	funcPtr := flag.String("function", "all", "Specifies which functions to run. Accepts 'all' wildcard")
	sleepPtr := flag.Duration("sleep", time.Minute*10, `A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`)

	flag.Parse()

	mysqlDbConn := config.MysqlDbConn()
	defer mysqlDbConn.Close()

	exchangeSlaveMysqlDbConn := config.ExchangeSlaveMysqlDbConn()
	defer exchangeSlaveMysqlDbConn.Close()

	healthCheckRepo := mysql.NewMysqlHealthCheckRepository(mysqlDbConn)
	systemConfigRepo := mysql.NewMysqlSystemConfigRepository(mysqlDbConn)

	coldbalanceRepo := mysql.NewMysqlColdBalanceRepository(mysqlDbConn)
	hotLimitRepo := exchange.NewExchangeHotLimitRepository()
	rpcMethodRepo := mysql.NewMysqlRpcMethodRepository(mysqlDbConn)
	rpcRequestRepo := mysql.NewMysqlRpcRequestRepository(mysqlDbConn)
	rpcResponseRepo := mysql.NewMysqlRpcResponseRepository(mysqlDbConn)
	userRepo := mysql.NewMysqlUserRepository(mysqlDbConn)
	withdrawRepo := mysql.NewMysqlWithdrawRepository(exchangeSlaveMysqlDbConn)

	marketRepo := exchange.NewExchangeMarketRepository()
	userBalanceRepo := mysql.NewMysqlUserBalanceRepository(exchangeSlaveMysqlDbConn)

	coldWalletService := hcw.NewColdWalletService(coldbalanceRepo)
	ModuleServices := modules.NewModuleServices(healthCheckRepo, systemConfigRepo, rpcMethodRepo, rpcRequestRepo, rpcResponseRepo)
	MarketService := h.NewMarketService(marketRepo)
	walletService := hw.NewWalletService(ModuleServices, coldWalletService, MarketService, withdrawRepo, hotLimitRepo, userBalanceRepo)

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Println("Running " + *funcPtr + " with " + sleepPtr.String() + " sleep")

	if *funcPtr == "all" || *funcPtr == "checkbalance" {
		logger.Log("Initializing checkbalance service...")
		checkBalanceService := hc.NewCheckBalanceService(walletService, coldWalletService, MarketService, ModuleServices, hotLimitRepo, userRepo)

		go func() {
			for {
				logger.Log("- Running checkbalance ...")
				checkBalanceService.CheckBalanceHandler()

				logger.Log("- Finished running checkbalance, sleeping for " + sleepPtr.String() + " ...")
				time.Sleep(*sleepPtr)
			}
		}()
	}

	if *funcPtr == "all" || *funcPtr == "healthcheck" {
		logger.Log("Initializing healthcheck service...")
		healthCheckService := hc.NewHealthCheckService(ModuleServices, healthCheckRepo, systemConfigRepo)
		
		go func() {
			for {
				logger.Log("- Running healthcheck ...")
				healthCheckService.HealthCheckHandler()

				logger.Log("- Finished running healthcheck, sleeping for " + sleepPtr.String() + " ...")
				time.Sleep(*sleepPtr)
			}
		}()
	}

	wg.Wait()
}
