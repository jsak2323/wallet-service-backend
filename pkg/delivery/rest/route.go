package rest

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/deposit"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/role"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/user"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/walletrpc"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/withdraw"
	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	routes      *mux.Router
	svc         service.Service
	permission  permission.PermissionHandler
	role        role.RoleHandler
	user        user.UserHandler
	currency    currency.CurrencyHandler
	rpcConfig   rpcconfig.RpcConfigHandler
	rpcMethod   rpcmethod.RpcMethodHandler
	rpcRequest  rpcrequest.RpcRequestHandler
	rpcresponse rpcresponse.RpcResponseHandler
	deposit     deposit.DepositHandler
	fireblocks  fireblocks.FireblocksHandler
	withdraw    withdraw.WithdrawHandler
	walletrpc   walletrpc.WalletRpcHandler
}

// New ...
func New(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes:      routes,
		svc:         svc,
		permission:  permission.NewPermissionHandler(routes, svc),
		role:        role.NewRoleHandler(routes, svc),
		user:        user.NewUserHandler(routes, svc),
		currency:    currency.NewCurrencyHandler(routes, svc),
		rpcConfig:   rpcconfig.NewRpcConfigHandler(routes, svc),
		rpcMethod:   rpcmethod.NewRpcMethodHandler(routes, svc),
		rpcRequest:  rpcrequest.NewRpcRequestHandler(routes, svc),
		rpcresponse: rpcresponse.NewRpcResponseHandler(routes, svc),
		deposit:     deposit.NewDepositHandler(routes, svc),
		fireblocks:  fireblocks.NewFireblocksHandler(routes, svc),
		withdraw:    withdraw.NewWithdrawHandler(routes, svc),
		walletrpc:   walletrpc.NewWalletRpcHandler(routes, svc),
	}
}

func (re *Rest) Route() {
	re.routes.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
	})

	re.routes.HandleFunc("/login", re.user.LoginHandler).Methods(http.MethodPost).Name("login")
	re.routes.HandleFunc("/logout", re.user.LogoutHandler).Methods(http.MethodPost)

	// -- User management
	user := re.routes.PathPrefix("/user").Subrouter()
	userHandler := re.user
	user.HandleFunc("/list", userHandler.ListUserHandler).Methods(http.MethodGet).Name("listusers")
	user.HandleFunc("", userHandler.CreateUserHandler).Methods(http.MethodPost).Name("createuser")
	user.HandleFunc("", userHandler.UpdateUserHandler).Methods(http.MethodPut).Name("updateuser")
	user.HandleFunc("/deactivate/{id}", userHandler.DeactivateUserHandler).Methods(http.MethodPost).Name("deactivateuser")
	user.HandleFunc("/activate/{id}", userHandler.ActivateUserHandler).Methods(http.MethodPost).Name("activateuser")
	user.HandleFunc("/role", userHandler.AddUserRolesHandler).Methods(http.MethodPost).Name("createuserrole")
	user.HandleFunc("/{user_id}/role/{role_id}", userHandler.DeleteUserRoleHandler).Methods(http.MethodDelete).Name("deleteuserrole")

	// -- Role management
	role := re.routes.PathPrefix("/role").Subrouter()
	roleHandler := re.role
	role.HandleFunc("/role/list", roleHandler.ListRoleHandler).Methods(http.MethodGet).Name("listroles")
	role.HandleFunc("/role", roleHandler.CreateRoleHandler).Methods(http.MethodPost).Name("createrole")
	role.HandleFunc("/role", roleHandler.UpdateRoleHandler).Methods(http.MethodPut).Name("updaterole")
	role.HandleFunc("/role/{id}", roleHandler.DeleteRoleHandler).Methods(http.MethodDelete).Name("deleterole")
	role.HandleFunc("/role/permission", roleHandler.CreateRolePermissionHandler).Methods(http.MethodPost).Name("createrolepermission")
	role.HandleFunc("/role/{role_id}/permission/{permission_id}", roleHandler.DeleteRolePermissionHandler).Methods(http.MethodDelete).Name("deleterolepermission")

	// -- Permission management
	permission := re.routes.PathPrefix("/permission").Subrouter()
	permissionHandler := re.permission
	permission.HandleFunc("/list", permissionHandler.ListPermissionHandler).Methods(http.MethodGet).Name("listpermissions")
	permission.HandleFunc("/permission", permissionHandler.CreatePermissionHandler).Methods(http.MethodPost).Name("createpermission")
	permission.HandleFunc("", permissionHandler.UpdatePermissionHandler).Methods(http.MethodPut).Name("updatepermission")
	permission.HandleFunc("/permission/{id}", permissionHandler.DeletePermissionHandler).Methods(http.MethodDelete).Name("deletepermission")

	// -- Currency Config management
	currency := re.routes.PathPrefix("/currency").Subrouter()
	currencyHandler := re.currency
	currency.HandleFunc("/list", currencyHandler.ListHandler).Methods(http.MethodGet).Name("listcurrency")
	currency.HandleFunc("", currencyHandler.CreateHandler).Methods(http.MethodPost).Name("createcurrency")
	currency.HandleFunc("", currencyHandler.UpdateHandler).Methods(http.MethodPut).Name("updatecurrency")
	currency.HandleFunc("/rpcconfig", currencyHandler.CreateRpcHandler).Methods(http.MethodPost).Name("createcurrencyrpc")
	currency.HandleFunc("/{currency_id}/rpcconfig/{rpc_id}", currencyHandler.DeleteRpcHandler).Methods(http.MethodDelete).Name("deletecurrencyrpc")
	currency.HandleFunc("/deactivate/{id}", currencyHandler.DeactivateHandler).Methods(http.MethodPost).Name("deactivatecurrency")
	currency.HandleFunc("/activate/{id}", currencyHandler.ActivateHandler).Methods(http.MethodPost).Name("activatecurrency")

	// -- Rpc Config management
	rpcConfig := re.routes.PathPrefix("/rpcconfig").Subrouter()
	rpcConfigHandler := re.rpcConfig
	rpcConfig.HandleFunc("/list", rpcConfigHandler.ListHandler).Methods(http.MethodGet).Name("listrpcconfig")
	rpcConfig.HandleFunc("/id/{id}", rpcConfigHandler.GetByIdHandler).Methods(http.MethodGet).Name("getrpcconfigbyid")
	rpcConfig.HandleFunc("", rpcConfigHandler.CreateHandler).Methods(http.MethodPost).Name("createrpcconfig")
	rpcConfig.HandleFunc("", rpcConfigHandler.UpdateHandler).Methods(http.MethodPut).Name("updaterpcconfig")
	rpcConfig.HandleFunc("/deactivate/{id}", rpcConfigHandler.DeactivateHandler).Methods(http.MethodPost).Name("deactivaterpcconfig")
	rpcConfig.HandleFunc("/rpcmethod", rpcConfigHandler.CreateRpcMethodHandler).Methods(http.MethodPost).Name("createrpcconfigrpcmethod")
	rpcConfig.HandleFunc("/{rpcconfig_id}/rpcmethod/{rpcmethod_id}", rpcConfigHandler.DeleteRpcMethodHandler).Methods(http.MethodDelete).Name("deleterpcconfigrpcmethod")
	rpcConfig.HandleFunc("/activate/{id}", rpcConfigHandler.ActivateHandler).Methods(http.MethodPost).Name("activaterpcconfig")

	rpcMethod := re.routes.PathPrefix("/rpcmethod").Subrouter()
	rpcMethodHandler := re.rpcMethod
	rpcMethod.HandleFunc("/rpcmethod/list", rpcMethodHandler.ListHandler).Methods(http.MethodGet).Name("listrpcmethod")
	rpcMethod.HandleFunc("/rpcmethod/rpcconfig/{rpc_config_id}", rpcMethodHandler.GetByRpcConfigIdHandler).Methods(http.MethodGet).Name("rpcmethodbyrpcconfig")
	rpcMethod.HandleFunc("/rpcmethod", rpcMethodHandler.CreateHandler).Methods(http.MethodPost).Name("createrpcmethod")
	rpcMethod.HandleFunc("/rpcmethod", rpcMethodHandler.UpdateHandler).Methods(http.MethodPut).Name("updaterpcmethod")
	rpcMethod.HandleFunc("/rpcmethod/{id}/rpcconfig/{rpc_config_id}", rpcMethodHandler.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcmethod")

	rpcRequest := re.routes.PathPrefix("/rpcrequest").Subrouter()
	rpcRequestHandler := re.rpcRequest
	rpcRequest.HandleFunc("/rpcmethod/{rpc_method_id}", rpcRequestHandler.GetByRpcMethodIdHandler).Methods(http.MethodGet).Name("rpcrequestbyrpcmethod")
	rpcRequest.HandleFunc("", rpcRequestHandler.CreateHandler).Methods(http.MethodPost).Name("createrpcrequest")
	rpcRequest.HandleFunc("", rpcRequestHandler.UpdateHandler).Methods(http.MethodPut).Name("updaterpcrequest")
	rpcRequest.HandleFunc("/{id}/rpcmethod/{rpc_method_id}", rpcRequestHandler.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcrequest")

	rpcResponse := re.routes.PathPrefix("/rpcresponse").Subrouter()
	rpcResponseHandler := re.rpcresponse
	rpcResponse.HandleFunc("/rpcmethod/{rpc_method_id}", rpcResponseHandler.GetByRpcMethodIdHandler).Methods(http.MethodGet).Name("rpcresponsebyrpcmethod")
	rpcResponse.HandleFunc("", rpcResponseHandler.CreateHandler).Methods(http.MethodPost).Name("createrpcresponse")
	rpcResponse.HandleFunc("", rpcResponseHandler.UpdateHandler).Methods(http.MethodPut).Name("updaterpcresponse")
	rpcResponse.HandleFunc("/{id}/rpcmethod/{rpc_method_id}", rpcResponseHandler.DeleteHandler).Methods(http.MethodDelete).Name("deleterpcresponse")

	deposit := re.routes.PathPrefix("/deposit").Subrouter()
	depositHandler := re.deposit
	deposit.HandleFunc("/deposit", depositHandler.ListHandler).Methods(http.MethodGet).Name("listdeposits")

	fireblocks := re.routes.PathPrefix("/fireblocks").Subrouter()
	fireblocksHandler := re.fireblocks
	fireblocks.HandleFunc("/tx_sign_request", fireblocksHandler.CallbackHandler).Methods(http.MethodPost).Name("fireblockscallback")

	withdraw := re.routes.PathPrefix("/withdraw").Subrouter()
	withdrawHandler := re.deposit
	withdraw.HandleFunc("", withdrawHandler.ListHandler).Methods(http.MethodGet).Name("listwithdraws")

	// command
	walletRpcHandler := re.walletrpc
	listTransactions := re.routes.PathPrefix("/listtransactions").Subrouter()
	listTransactions.HandleFunc("", walletRpcHandler.ListTransactionsHandler).Methods(http.MethodGet)
	listTransactions.HandleFunc("/{limit}", walletRpcHandler.ListTransactionsHandler).Methods(http.MethodGet)
	re.routes.HandleFunc("/{symbol}/listtransactions", walletRpcHandler.ListTransactionsHandler).Methods(http.MethodGet)
	re.routes.HandleFunc("/{symbol}/listtransactions/{limit}", walletRpcHandler.ListTransactionsHandler).Methods(http.MethodGet)

	sendToAddress := re.routes.PathPrefix("/sendtoaddress").Subrouter()
	sendToAddress.HandleFunc("", walletRpcHandler.SendToAddressHandler).Methods(http.MethodPost).Name("sendhotwallet")

	re.routes.HandleFunc("/{symbol}/addresstype/{address}", walletRpcHandler.AddressTypeHandler).Methods(http.MethodGet)

	nodes := re.routes.PathPrefix("/nodes").Subrouter()
	nodes.HandleFunc("/getbalance", walletRpcHandler.GetBalanceHandler).Methods(http.MethodGet)
	nodes.HandleFunc("/{symbol}/getbalance", walletRpcHandler.GetBalanceHandler).Methods(http.MethodGet)

	wallet := re.routes.PathPrefix("/wallet").Subrouter()
	wallet.HandleFunc("/getbalance", walletRpcHandler.GetBalanceHandler).Methods(http.MethodGet).Name("listbalances")
	wallet.HandleFunc("/{currency_id}/getbalance", walletRpcHandler.GetBalanceHandler).Methods(http.MethodGet).Name("balancebycurrencyid")

	re.routes.HandleFunc("/getblockcount", walletRpcHandler.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcount")
	re.routes.HandleFunc("/{symbol}/getblockcount", walletRpcHandler.GetBlockCountHandler).Methods(http.MethodGet).Name("getblockcountbysymbol")

	re.routes.HandleFunc("/gethealthcheck", walletRpcHandler.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheck")
	re.routes.HandleFunc("/{symbol}/gethealthcheck", walletRpcHandler.GetHealthCheckHandler).Methods(http.MethodGet).Name("gethealthcheckbysymbol")

	re.routes.HandleFunc("/{symbol}/getnewaddress", walletRpcHandler.GetNewAddressHandler).Methods(http.MethodGet)
	re.routes.HandleFunc("/{symbol}/getnewaddress/{type}", walletRpcHandler.GetNewAddressHandler).Methods(http.MethodGet)

	re.routes.HandleFunc("/log/{symbol}/{rpcconfigtype}/{date}", walletRpcHandler.GetLogHandler).Methods(http.MethodGet).Name("getlog")

	// // r.HandleFunc("/userwallet/getbalance", userWalletService.GetBalanceHandler).Methods(http.MethodGet)

}
