package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	hd "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/deposit"
	hp "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/permission"
	hr "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	hrc "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	hrm "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcmethod"
	hrrq "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcrequest"
	hrrs "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcresponse"
	hu "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	huw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/user"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	authm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/auth"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

func setRoutes(r *mux.Router, mysqlRepos mysql.MysqlRepositories, exchangeApiRepos exchange.APIRepositories, validator *util.CustomValidator) {
	// -- Auth
	userService := hu.NewUserService(mysqlRepos.User, mysqlRepos.Role, mysqlRepos.UserRole, mysqlRepos.Permission, *validator)
	r.HandleFunc("/login", userService.LoginHandler).Methods(http.MethodPost).Name("login")
	r.HandleFunc("/logout", userService.LogoutHandler).Methods(http.MethodPost)

	// MODULE SERVICES
	ModuleServices := modules.NewModuleServices(mysqlRepos.HealthCheck, mysqlRepos.SystemConfig, mysqlRepos.RpcMethod, mysqlRepos.RpcRequest, mysqlRepos.RpcResponse)
	MarketService := h.NewMarketService(exchangeApiRepos.Market)

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
	roleService := hr.NewRoleService(mysqlRepos.Role, mysqlRepos.Permission, mysqlRepos.RolePermission, mysqlRepos.UserRole, *validator)
	r.HandleFunc("/role/list", roleService.ListRoleHandler).Methods(http.MethodGet).Name("listroles")
	r.HandleFunc("/role", roleService.CreateRoleHandler).Methods(http.MethodPost).Name("createrole")
	r.HandleFunc("/role", roleService.UpdateRoleHandler).Methods(http.MethodPut).Name("updaterole")
	r.HandleFunc("/role/{id}", roleService.DeleteRoleHandler).Methods(http.MethodDelete).Name("deleterole")
	r.HandleFunc("/role/permission", roleService.CreatePermissionHandler).Methods(http.MethodPost).Name("createrolepermission")
	r.HandleFunc("/role/{role_id}/permission/{permission_id}", roleService.DeletePermissionHandler).Methods(http.MethodDelete).Name("deleterolepermission")

	// -- Permission management
	permissionService := hp.NewPermissionService(mysqlRepos.Permission, mysqlRepos.RolePermission, *validator)
	r.HandleFunc("/permission/list", permissionService.ListPermissionHandler).Methods(http.MethodGet).Name("listpermissions")
	r.HandleFunc("/permission", permissionService.CreatePermissionHandler).Methods(http.MethodPost).Name("createpermission")
	r.HandleFunc("/permission", permissionService.UpdatePermissionHandler).Methods(http.MethodPut).Name("updatepermission")
	r.HandleFunc("/permission/{id}", permissionService.DeletePermissionHandler).Methods(http.MethodDelete).Name("deletepermission")

	// -- GET getblockcount
	getBlockCountService := h.NewGetBlockCountService(ModuleServices, mysqlRepos.SystemConfig)
	r.HandleFunc("/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcount")
	r.HandleFunc("/{symbol}/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcountbysymbol")

	// -- GET gethealthcheck
	getHealthCheckService := h.NewGetHealthCheckService(ModuleServices, mysqlRepos.HealthCheck, mysqlRepos.SystemConfig)
	r.HandleFunc("/gethealthcheck", getHealthCheckService.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheck")
	r.HandleFunc("/{symbol}/gethealthcheck", getHealthCheckService.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheckbysymbol")

	// -- GET log
	getLogService := h.NewGetLogService(ModuleServices)
	r.HandleFunc("/log/{symbol}/{rpcconfigtype}/{date}", getLogService.GetLogHandler).Methods(http.MethodGet).Name("getlog")

	// -- Currency Config management
	currencyConfigService := hc.NewCurrencyConfigService(mysqlRepos.CurrencyConfig, mysqlRepos.CurrencyRpc, mysqlRepos.RpcConfig, *validator)
	r.HandleFunc("/currency/list", currencyConfigService.ListHandler).Methods(http.MethodGet).Name("listcurrency")
	r.HandleFunc("/currency", currencyConfigService.CreateHandler).Methods(http.MethodPost).Name("createcurrency")
	r.HandleFunc("/currency", currencyConfigService.UpdateHandler).Methods(http.MethodPut).Name("updatecurrency")
	r.HandleFunc("/currency/rpcconfig", currencyConfigService.CreateRpcHandler).Methods(http.MethodPost).Name("createcurrencyrpc")
	r.HandleFunc("/currency/{currency_id}/rpcconfig/{rpc_id}", currencyConfigService.DeleteRpcHandler).Methods(http.MethodDelete).Name("deletecurrencyrpc")
	r.HandleFunc("/currency/deactivate/{id}", currencyConfigService.DeactivateHandler).Methods(http.MethodPost).Name("deactivatecurrency")
	r.HandleFunc("/currency/activate/{id}", currencyConfigService.ActivateHandler).Methods(http.MethodPost).Name("activatecurrency")
	// -- GET getbalance
	getBalanceService := h.NewGetBalanceService(ModuleServices)
	r.HandleFunc("/nodes/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
	r.HandleFunc("/nodes/{symbol}/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)

	// -- Rpc Config management
	rpcConfigService := hrc.NewRpcConfigService(mysqlRepos.RpcConfig, mysqlRepos.RpcConfigRpcMethod, *validator)
	r.HandleFunc("/rpcconfig/list", rpcConfigService.ListHandler).Methods(http.MethodGet).Name("listrpcconfig")
	r.HandleFunc("/rpcconfig/id/{id}", rpcConfigService.GetByIdHandler).Methods(http.MethodGet).Name("getrpcconfigbyid")
	r.HandleFunc("/rpcconfig", rpcConfigService.CreateHandler).Methods(http.MethodPost).Name("createrpcconfig")
	r.HandleFunc("/rpcconfig", rpcConfigService.UpdateHandler).Methods(http.MethodPut).Name("updaterpcconfig")
	r.HandleFunc("/rpcconfig/deactivate/{id}", rpcConfigService.DeactivateHandler).Methods(http.MethodPost).Name("deactivaterpcconfig")
	r.HandleFunc("/rpcconfig/rpcmethod", rpcConfigService.CreateRpcMethodHandler).Methods(http.MethodPost).Name("createrpcconfigrpcmethod")
	r.HandleFunc("/rpcconfig/{rpcconfig_id}/rpcmethod/{rpcmethod_id}", rpcConfigService.DeleteRpcMethodHandler).Methods(http.MethodDelete).Name("deleterpcconfigrpcmethod")
	r.HandleFunc("/rpcconfig/activate/{id}", rpcConfigService.ActivateHandler).Methods(http.MethodPost).Name("activaterpcconfig")

	rpcMethodService := hrm.NewRpcMethodService(mysqlRepos.RpcMethod, mysqlRepos.RpcConfigRpcMethod, *validator)
	r.HandleFunc("/rpcmethod/list", rpcMethodService.ListHandler).Methods(http.MethodGet).Name("listrpcmethod")
	r.HandleFunc("/rpcmethod/rpcconfig/{rpc_config_id}", rpcMethodService.GetByRpcConfigIdHandler).Methods(http.MethodGet).Name("rpcmethodbyrpcconfig")
	r.HandleFunc("/rpcmethod", rpcMethodService.CreateHandler).Methods(http.MethodPost).Name("createrpcmethod")
	r.HandleFunc("/rpcmethod", rpcMethodService.UpdateHandler).Methods(http.MethodPut).Name("updaterpcmethod")
	r.HandleFunc("/rpcmethod/{id}/rpcconfig/{rpc_config_id}", rpcMethodService.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcmethod")

	rpcRequestService := hrrq.NewRpcRequestService(mysqlRepos.RpcRequest, *validator)
	r.HandleFunc("/rpcrequest/rpcmethod/{rpc_method_id}", rpcRequestService.GetByRpcMethodIdHandler).Methods(http.MethodGet).Name("rpcrequestbyrpcmethod")
	r.HandleFunc("/rpcrequest", rpcRequestService.CreateHandler).Methods(http.MethodPost).Name("createrpcrequest")
	r.HandleFunc("/rpcrequest", rpcRequestService.UpdateHandler).Methods(http.MethodPut).Name("updaterpcrequest")
	r.HandleFunc("/rpcrequest/{id}/rpcmethod/{rpc_method_id}", rpcRequestService.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcrequest")

	rpcResponseService := hrrs.NewRpcResponseService(mysqlRepos.RpcResponse, *validator)
	r.HandleFunc("/rpcresponse/rpcmethod/{rpc_method_id}", rpcResponseService.GetByRpcMethodIdHandler).Methods(http.MethodGet).Name("rpcresponsebyrpcmethod")
	r.HandleFunc("/rpcresponse", rpcResponseService.CreateHandler).Methods(http.MethodPost).Name("createrpcresponse")
	r.HandleFunc("/rpcresponse", rpcResponseService.UpdateHandler).Methods(http.MethodPut).Name("updaterpcresponse")
	r.HandleFunc("/rpcresponse/{id}/rpcmethod/{rpc_method_id}", rpcResponseService.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcresponse")

	// -- Cold Wallet management
	coldWalletService := hcw.NewColdWalletService(mysqlRepos.ColdBalance, *validator)
	r.HandleFunc("/coldwallet", coldWalletService.CreateHandler).Methods(http.MethodPost).Name("createcoldwallet")
	// TODO get by currency
	r.HandleFunc("/coldwallet", coldWalletService.ListHandler).Methods(http.MethodGet).Name("listcoldwallet")
	r.HandleFunc("/coldwallet/sendtohot", coldWalletService.SendToHotHandler).Methods(http.MethodPost).Name("sendcoldwallet")
	r.HandleFunc("/coldwallet", coldWalletService.UpdateHandler).Methods(http.MethodPut).Name("updatecoldwallet")
	r.HandleFunc("/coldwallet/updatebalance", coldWalletService.UpdateBalanceHandler).Methods(http.MethodPut).Name("updatecoldwalletbalance")
	r.HandleFunc("/coldwallet/deactivate/{id}", coldWalletService.DeactivateHandler).Methods(http.MethodPost).Name("deactivatecoldwallet")
	r.HandleFunc("/coldwallet/activate/{id}", coldWalletService.ActivateHandler).Methods(http.MethodPost).Name("activatecoldwallet")

	userWalletService := huw.NewUserWalletService(mysqlRepos.UserBalance)
	r.HandleFunc("/userwallet/getbalance", userWalletService.GetBalanceHandler).Methods(http.MethodGet)

	walletService := hw.NewWalletService(ModuleServices, coldWalletService, MarketService, mysqlRepos.Withdraw, exchangeApiRepos.HotLimit, mysqlRepos.UserBalance)
	r.HandleFunc("/wallet/getbalance", walletService.GetBalanceHandler).Methods(http.MethodGet).Name("listbalances")
	r.HandleFunc("/wallet/{currency_id}/getbalance", walletService.GetBalanceHandler).Methods(http.MethodGet).Name("balancebycurrencyid")

	depositService := hd.NewDepositService(mysqlRepos.Deposit)
	r.HandleFunc("/deposit", depositService.ListHandler).Methods(http.MethodGet).Name("listdeposits")

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
	systemConfigService := h.NewSystemConfigService(mysqlRepos.SystemConfig)
	r.HandleFunc("/systemconfig/maintenancelist/{action}/{value}", systemConfigService.MaintenanceListHandler).Methods(http.MethodPut).Name("updatemaintlist")
	/*
	   curl example:
	   curl --request PUT localhost:3000/systemconfig/maintenancelist/add/BTC
	*/

	fireblocksService := h.NewFireblocksService()
	r.HandleFunc("/fireblocks/tx_sign_request", fireblocksService.CallbackHandler).Methods(http.MethodPost).Name("fireblockscallback")

	auth := authm.NewAuthMiddleware(mysqlRepos.Role, mysqlRepos.Permission, mysqlRepos.RolePermission)
	r.Use(auth.Authenticate)
	r.Use(auth.Authorize)
}
