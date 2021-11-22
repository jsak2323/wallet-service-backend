package cron

import (
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func Run(function string, sleep time.Duration, mysqlRepos mysql.MysqlRepositories, exchangeApiRepos exchange.APIRepositories) {
	switch function {
	case "checkbalance":
		runCheckBalance(sleep, mysqlRepos, exchangeApiRepos)
	case "healthcheck":
		runHealthCheck(sleep, mysqlRepos)
	case "default":
		fmt.Println("function not specified")
	}
}

func countDownSleep(function string, delay int) {
	ticker := time.Tick(time.Second)

	for i := delay; i >= 0; i-- {
		<-ticker
		fmt.Printf("\r - Next "+function+" execution in %d seconds ...", i)
	}
}
