package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/cmd/server/cron"
	"github.com/btcid/wallet-services-backend-go/cmd/server/http"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
	"github.com/go-playground/validator"
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	localMysqlDbConn := config.MysqlDbConn()

	exchangeSlaveMysqlDbConn := config.ExchangeSlaveMysqlDbConn()

	go func() {
		sigchan := make(chan os.Signal)

		signal.Notify(sigchan, os.Interrupt)

		<-sigchan

		exchangeSlaveMysqlDbConn.Close()
		fmt.Println("database exchange has been closed")

		localMysqlDbConn.Close()
		fmt.Println("local database has been closed")

		os.Exit(0)
	}()

	mysqlRepos := mysql.NewMysqlRepositories(localMysqlDbConn, exchangeSlaveMysqlDbConn)
	exchangeApiRepos := exchange.NewAPIRepositories()

	validator := &util.CustomValidator{Validator: validator.New()}

	switch *appPtr {
	case "http":
		http.Run(mysqlRepos, exchangeApiRepos, validator)
	case "cron":
		cron.Run(*funcPtr, *sleepPtr, mysqlRepos, exchangeApiRepos, validator)
	}
}
