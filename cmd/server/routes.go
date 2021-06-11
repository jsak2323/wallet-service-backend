package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	c "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	authm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/auth"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/go-redis/redis/v8"
)

func SetRoutes(r *mux.Router, mysqlDbConn *sql.DB, redis *redis.Client) {
	// REPOSITORIES
	userRepo := mysql.NewMysqlUserRepository(mysqlDbConn)
	userRoleRepo := mysql.NewMysqlUserRoleRepository(mysqlDbConn)
	roleRepo := mysql.NewMysqlRoleRepository(mysqlDbConn)
	permissionRepo := mysql.NewMysqlPermissionRepository(mysqlDbConn)
	rolePermissionRepo := mysql.NewMysqlRolePermissionRepository(mysqlDbConn)

	healthCheckRepo := mysql.NewMysqlHealthCheckRepository(mysqlDbConn)
	systemConfigRepo := mysql.NewMysqlSystemConfigRepository(mysqlDbConn)

	// -- User
	userService := h.NewUserService(userRepo, userRoleRepo, roleRepo, redis)
	r.HandleFunc("/register", userService.RegisterHandler).Methods(http.MethodPost).Name("register")
	r.HandleFunc("/login", userService.LoginHandler).Methods(http.MethodPost).Name("login")
	r.HandleFunc("/logout", userService.LogoutHandler).Methods(http.MethodPost)

	// MODULE SERVICES
	ModuleServices := modules.NewModuleServices(healthCheckRepo, systemConfigRepo)

	// API ROUTES

	// -- GET getblockcount
	getBlockCountService := h.NewGetBlockCountService(ModuleServices, systemConfigRepo)
	r.HandleFunc("/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcount")
	r.HandleFunc("/{symbol}/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcountbysymbol")

	// -- GET gethealthcheck
	getHealthCheckService := h.NewGetHealthCheckService(ModuleServices, healthCheckRepo, systemConfigRepo)
	r.HandleFunc("/gethealthcheck", getHealthCheckService.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheck")
	r.HandleFunc("/{symbol}/gethealthcheck", getHealthCheckService.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheckbysymbol")

	// -- GET log
	getLogService := h.NewGetLogService(ModuleServices)
	r.HandleFunc("/log/{symbol}/{rpcconfigtype}/{date}", getLogService.GetLogHandler).Methods(http.MethodGet).Name("getlog")

	// -- GET getbalance (disabled)
	/*
	   getBalanceService := h.NewGetBalanceService(ModuleServices)
	   r.HandleFunc("/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
	   r.HandleFunc("/{symbol}/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
	*/

	// -- GET listtransactions (disabled)
	/*
	   listTransactionsService := h.NewListTransactionsService(ModuleServices)
	   r.HandleFunc("/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
	   r.HandleFunc("/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
	   r.HandleFunc("/{symbol}/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
	   r.HandleFunc("/{symbol}/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
	*/

	// -- GET getnewaddress (disabled)
	/*
	   getNewAddressService := h.NewGetNewAddressService(ModuleServices)
	   r.HandleFunc("/{symbol}/getnewaddress", getNewAddressService.GetNewAddressHandler).Methods(http.MethodGet)
	   r.HandleFunc("/{symbol}/getnewaddress/{type}", getNewAddressService.GetNewAddressHandler).Methods(http.MethodGet)
	*/

	// -- GET addresstype (disabled)
	/*
	   addressTypeService := h.NewAddressTypeService(ModuleServices)
	   r.HandleFunc("/{symbol}/addresstype/{address}", addressTypeService.AddressTypeHandler).Methods(http.MethodGet)
	*/

	// -- POST sendtoaddress (disabled)
	/*
	   sendToAddressService := h.NewSendToAddressService(ModuleServices)
	   r.HandleFunc("/sendtoaddress", sendToAddressService.SendToAddressHandler).Methods(http.MethodPost)
	*/
	/*
	   curl example:
	   curl --header "Content-Type: application/json" --request POST --data '{"symbol":"btc", "amount":"0.001", "address":"2MtU6EMx37AYrCNj1RcRr6bw66QqHYw4D4R"}' localhost:3000/sendtoaddress | jq
	*/

	// -- PUT systemconfig
	systemConfigService := h.NewSystemConfigService(systemConfigRepo)
	r.HandleFunc("/systemconfig/maintenancelist/{action}/{value}", systemConfigService.MaintenanceListHandler).Methods(http.MethodPut).Name("updatemaintlist")
	/*
	   curl example:
	   curl --request PUT localhost:3000/systemconfig/maintenancelist/add/BTC
	*/

	// CRON ROUTES

	// -- GET healthcheck
	healthCheckService := c.NewHealthCheckService(ModuleServices, healthCheckRepo, systemConfigRepo)
	r.HandleFunc("/cron/healthcheck", healthCheckService.HealthCheckHandler).Methods(http.MethodGet).Name("cronhealthcheck")

	auth := authm.NewAuthMiddleware(roleRepo, permissionRepo, rolePermissionRepo, redis)
	r.Use(auth.Authenticate)
	r.Use(auth.Authorize)
}
