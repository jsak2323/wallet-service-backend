package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/cmd/server/cron"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	authm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/auth"
	logm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func serv() {
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

	switch *appPtr {
	case "http":
		// http.Run(mysqlRepos, exchangeApiRepos, validator)
		routes := mux.NewRouter()

		rest.New(
			routes,
			service,
		).Route()

		corsOpts := cors.New(cors.Options{
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
				http.MethodOptions,
				http.MethodHead,
			},
			AllowedHeaders: []string{
				"*",
			},
		})

		routes.Use(logm.LogMiddleware)

		auth := authm.NewAuthMiddleware(mysqlRepos.Role, mysqlRepos.Permission, mysqlRepos.RolePermission)
		routes.Use(auth.Authenticate)
		routes.Use(auth.Authorize)

		server := &http.Server{
			Handler:      corsOpts.Handler(routes),
			Addr:         ":" + config.CONF.Port,
			WriteTimeout: 120 * time.Second,
			ReadTimeout:  120 * time.Second,
		}

		fmt.Println()
		fmt.Println("Running server on localhost:" + config.CONF.Port)
		fmt.Println("\n\n\n")

		log.Fatal(server.ListenAndServe())
	case "cron":
		cron.Run(*funcPtr, *sleepPtr, mysqlRepos, exchangeApiRepos, validator)
	}
}
