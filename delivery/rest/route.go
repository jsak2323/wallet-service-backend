package rest

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/service"
	"github.com/gorilla/mux"
)

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

// New ...
func New(
	svc service.Service,
	routes *mux.Router,
) *Rest {
	return &Rest{
		svc:    svc,
		routes: routes,
	}
}

func (re *Rest) Route() {
	re.routes.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
	})

	re.routes.HandleFunc("/login", re.LoginHandler).Methods(http.MethodPost).Name("login")
	re.routes.HandleFunc("/logout", re.LogoutHandler).Methods(http.MethodPost)

	// -- User management
	user := re.routes.PathPrefix("/user").Subrouter()
	user.HandleFunc("/list", re.ListUserHandler).Methods(http.MethodGet).Name("listusers")
	user.HandleFunc("", re.CreateUserHandler).Methods(http.MethodPost).Name("createuser")
	user.HandleFunc("", re.UpdateUserHandler).Methods(http.MethodPut).Name("updateuser")
	user.HandleFunc("/deactivate/{id}", re.DeactivateUserHandler).Methods(http.MethodPost).Name("deactivateuser")
	user.HandleFunc("/activate/{id}", re.ActivateUserHandler).Methods(http.MethodPost).Name("activateuser")
	user.HandleFunc("/role", re.AddRolesHandler).Methods(http.MethodPost).Name("createuserrole")
	user.HandleFunc("/{user_id}/role/{role_id}", re.DeleteRoleHandler).Methods(http.MethodDelete).Name("deleteuserrole")

	// -- Permission management
	permission := re.routes.PathPrefix("/permission").Subrouter()
	permission.HandleFunc("/list", re.ListPermissionHandler).Methods(http.MethodGet).Name("listpermissions")
	permission.HandleFunc("/permission", re.CreatePermissionHandler).Methods(http.MethodPost).Name("createpermission")
	permission.HandleFunc("", re.UpdatePermissionHandler).Methods(http.MethodPut).Name("updatepermission")
	permission.HandleFunc("/permission/{id}", re.DeletePermissionHandler).Methods(http.MethodDelete).Name("deletepermission")
}
