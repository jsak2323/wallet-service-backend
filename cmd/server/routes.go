package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	huw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/user"
	hu "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	hr "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	hp "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/permission"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	c "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/cron"
	authm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/auth"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

func SetRoutes(r *mux.Router, mysqlDbConn *sql.DB, exchangeSlaveMysqlDbConn *sql.DB) {
	// REPOSITORIES
	userRepo := mysql.NewMysqlUserRepository(mysqlDbConn)
	roleRepo := mysql.NewMysqlRoleRepository(mysqlDbConn)
	urRepo := mysql.NewMysqlUserRoleRepository(mysqlDbConn)
	permissionRepo := mysql.NewMysqlPermissionRepository(mysqlDbConn)
	rolePermissionRepo := mysql.NewMysqlRolePermissionRepository(mysqlDbConn)

	currencyConfigRepo := mysql.NewMysqlCurrencyConfigRepository(mysqlDbConn)
	healthCheckRepo := mysql.NewMysqlHealthCheckRepository(mysqlDbConn)
	systemConfigRepo := mysql.NewMysqlSystemConfigRepository(mysqlDbConn)

	coldbalanceRepo := mysql.NewMysqlColdBalanceRepository(mysqlDbConn)
	hotLimitRepo := exchange.NewExchangeHotLimitRepository()

	userBalanceRepo := mysql.NewMysqlUserBalanceRepository(exchangeSlaveMysqlDbConn)
	withdrawRepo := mysql.NewMysqlWithdrawRepository(exchangeSlaveMysqlDbConn)
	marketRepo := exchange.NewExchangeMarketRepository()

	// -- Auth
	userService := hu.NewUserService(userRepo, roleRepo, urRepo, permissionRepo)
	r.HandleFunc("/login", userService.LoginHandler).Methods(http.MethodPost).Name("login")
	r.HandleFunc("/logout", userService.LogoutHandler).Methods(http.MethodPost)

	// MODULE SERVICES
	ModuleServices := modules.NewModuleServices(healthCheckRepo, systemConfigRepo)
	MarketService := h.NewMarketService(marketRepo)

	// API ROUTES

	// -- User management
	r.HandleFunc("/user/list", userService.ListUserHandler).Methods(http.MethodGet).Name("listusers")
	r.HandleFunc("/user", userService.CreateUserHandler).Methods(http.MethodPost).Name("createuser")
	r.HandleFunc("/user", userService.UpdateUserHandler).Methods(http.MethodPut).Name("updateuser")
	r.HandleFunc("/user/deactivate/{id}", userService.DeactivateUserHandler).Methods(http.MethodPost).Name("deactivateuser")
	r.HandleFunc("/user/activate/{id}", userService.ActivateUserHandler).Methods(http.MethodPost).Name("activateuser")
	r.HandleFunc("/user/role", userService.AddRolesHandler).Methods(http.MethodPost).Name("createuserrole")
	r.HandleFunc("/user/{user_id}/role/{role_id}", userService.DeleteRoleHandler).Methods(http.MethodDelete).Name("deleteuserrole")

	// -- Role management
	roleService := hr.NewRoleService(roleRepo, permissionRepo, rolePermissionRepo, urRepo)
	r.HandleFunc("/role/list", roleService.ListRoleHandler).Methods(http.MethodGet).Name("listroles")
	r.HandleFunc("/role", roleService.CreateRoleHandler).Methods(http.MethodPost).Name("createrole")
	r.HandleFunc("/role", roleService.UpdateRoleHandler).Methods(http.MethodPut).Name("updaterole")
	r.HandleFunc("/role/{id}", roleService.DeleteRoleHandler).Methods(http.MethodDelete).Name("deleterole")
	r.HandleFunc("/role/permission", roleService.CreatePermissionHandler).Methods(http.MethodPost).Name("createrolepermission")
	r.HandleFunc("/role/{role_id}/permission/{permission_id}", roleService.DeletePermissionHandler).Methods(http.MethodDelete).Name("deleterolepermission")

	// -- Permission management
	permissionService := hp.NewPermissionService(permissionRepo, rolePermissionRepo)
	r.HandleFunc("/permission/list", permissionService.ListPermissionHandler).Methods(http.MethodGet).Name("listpermissions")
	r.HandleFunc("/permission", permissionService.CreatePermissionHandler).Methods(http.MethodPost).Name("createpermission")
	r.HandleFunc("/permission", permissionService.UpdatePermissionHandler).Methods(http.MethodPut).Name("updatepermission")
	r.HandleFunc("/permission/{id}", permissionService.DeletePermissionHandler).Methods(http.MethodDelete).Name("deletepermission")

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

	// -- GET currencyconfig
	currencyConfigService := hc.NewCurrencyConfigService(currencyConfigRepo)
	r.HandleFunc("/currency/list", currencyConfigService.GetCurrencyConfigHandler).Methods(http.MethodGet).Name("listcurrency")
	
	// -- GET getbalance
	getBalanceService := h.NewGetBalanceService(ModuleServices)
	r.HandleFunc("/nodes/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/nodes/{symbol}/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)

	coldWalletService := hcw.NewColdWalletService(coldbalanceRepo)
	r.HandleFunc("/coldwallet/getbalance", coldWalletService.GetBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/coldwallet/{symbol}/getbalance", coldWalletService.GetBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/coldwallet/sendtohot", coldWalletService.SendToHotHandler).Methods(http.MethodPost).Name("sendcoldwallet")
	r.HandleFunc("/coldwallet/update", coldWalletService.UpdateBalanceHandler).Methods(http.MethodPost).Name("updatecoldwallet")

	userWalletService := huw.NewUserWalletService(userBalanceRepo)
	r.HandleFunc("/userwallet/getbalance", userWalletService.GetBalanceHandler).Methods(http.MethodGet)

	WalletService := hw.NewWalletService(ModuleServices, *coldWalletService, *MarketService, withdrawRepo, hotLimitRepo, userBalanceRepo)

	r.HandleFunc("/wallet/getbalance", WalletService.GetBalanceHandler).Methods(http.MethodGet).Name("listbalances")
	r.HandleFunc("/wallet/{symbol}/getbalance", WalletService.GetBalanceHandler).Methods(http.MethodGet).Name("balancebysymbol")
	
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
	sendToAddressService := h.NewSendToAddressService(ModuleServices)
	r.HandleFunc("/sendtoaddress", sendToAddressService.SendToAddressHandler).Methods(http.MethodPost).Name("sendhotwallet")
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

	// -- GET checkbalance
	checkBalanceService := c.NewCheckBalanceService(WalletService, coldWalletService, MarketService, *ModuleServices, hotLimitRepo, userRepo)
	r.HandleFunc("/cron/checkbalance", checkBalanceService.CheckBalanceHandler).Methods(http.MethodGet)

	auth := authm.NewAuthMiddleware(roleRepo, permissionRepo, rolePermissionRepo)
	r.Use(auth.Authenticate)
	r.Use(auth.Authorize)
}
