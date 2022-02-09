package rest

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/role"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/delivery/rest/user"

	// "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	routes     *mux.Router
	svc        service.Service
	permission permission.PermissionHandler
	role       role.RoleHandler
	user       user.UserHandler
	currency   currency.CurrencyHandler
	rpcConfig  rpcconfig.RpcConfigHandler
	rpcMethod  rpcmethod.RpcMethodHandler
	rpcRequest rpcrequest.RpcRequestHandler
}

// New ...
func New(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes:     routes,
		svc:        svc,
		permission: permission.NewPermissionHandler(routes, svc),
		role:       role.NewRoleHandler(routes, svc),
		user:       user.NewUserHandler(routes, svc),
		currency:   currency.NewCurrencyHandler(routes, svc),
		rpcConfig:  rpcconfig.NewRpcConfigHandler(routes, svc),
		rpcMethod:  rpcmethod.NewRpcMethodHandler(routes, svc),
		rpcRequest: rpcrequest.NewRpcRequestHandler(routes, svc),
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
}
