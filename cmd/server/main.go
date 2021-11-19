package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/cmd/server/cron"
	"github.com/btcid/wallet-services-backend-go/cmd/server/http"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func main() {
	appPtr := flag.String("app", "http", "Specifies which app to run.")
	funcPtr := flag.String("function", "all", "Specifies which functions to run.")
	sleepPtr := flag.Duration("sleep", time.Minute*10, `A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(" - panic: ", err)
		}
	}()

	flag.Parse()

	localMysqlDbConn := config.MysqlDbConn()
	defer localMysqlDbConn.Close()

	exchangeSlaveMysqlDbConn := config.ExchangeSlaveMysqlDbConn()
	defer exchangeSlaveMysqlDbConn.Close()

	mysqlRepos := mysql.NewMysqlRepositories(localMysqlDbConn, exchangeSlaveMysqlDbConn)
	exchangeApiRepos := exchange.NewAPIRepositories()

	switch *appPtr {
	case "http":
		http.Run(mysqlRepos, exchangeApiRepos)
	case "cron":
		cron.Run(*funcPtr, *sleepPtr, mysqlRepos, exchangeApiRepos)
	}
}
