package cron

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
	"github.com/go-playground/validator"
)

// ./main -app cron -function [[function name]] -sleep [[ sleep duration ]]

func Run(args []string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(" - panic: ", err)
		}
	}()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	localMysqlDbConn := config.MysqlDbConn()
	exchangeSlaveMysqlDbConn := config.ExchangeSlaveMysqlDbConn()

	go func() {
		sigchan := make(chan os.Signal, 1)

		signal.Notify(sigchan, os.Interrupt)

		<-sigchan

		exchangeSlaveMysqlDbConn.Close()
		fmt.Println("exchange database has been closed")

		localMysqlDbConn.Close()
		fmt.Println("local database has been closed")

		os.Exit(0)
	}()

	mysqlRepos := mysql.NewMysqlRepositories(localMysqlDbConn, exchangeSlaveMysqlDbConn)
	exchangeApiRepos := exchange.NewAPIRepositories()

	validator := &util.CustomValidator{Validator: validator.New()}

	var service = service.New(*validator, mysqlRepos, exchangeApiRepos)

	_ = service

	// -app cron -function [[function name]] -sleep [[ sleep duration ]]

	if len(args) != 3 {
		fmt.Println(" - panic: invalid arguments, must be on 3 arguments")
		return
	}

	switch args[0] {
	case "checkbalance":
		service.Cron.CheckBalanceHandler()
		// runCheckBalance(sleep, mysqlRepos, exchangeApiRepos, validator)
	case "healthcheck":
		service.Cron.HealthCheckHandler()
		// runHealthCheck(sleep, mysqlRepos)
	case "updatedeposit":
		service.Cron.UpdateDeposit()
		// runUpdateDeposit(sleep, mysqlRepos)
	case "updatewithdraw":
		service.Cron.UpdateWithdraw()
		// runUpdateWithdraw(sleep, mysqlRepos, exchangeApiRepos)
	default:
		fmt.Println("function not specified")
		return
	}

	// cron.Run(*funcPtr, *sleepPtr, mysqlRepos, exchangeApiRepos, validator)
	// }
}

func countDownSleep(function string, delay int) {
	ticker := time.Tick(time.Second)

	for i := delay; i >= 0; i-- {
		<-ticker
		fmt.Printf("\r - Next "+function+" execution in %d seconds ...", i)
	}
}
