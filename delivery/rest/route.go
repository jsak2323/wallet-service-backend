package rest

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
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

	// -- Permission management
	permission := re.routes.PathPrefix("/permission").Subrouter()
	permission.HandleFunc("/list", re.ListPermissionHandler).Methods(http.MethodGet).Name("listpermissions")
	// r.HandleFunc("/permission", permissionService.CreatePermissionHandler).Methods(http.MethodPost).Name("createpermission")
	permission.HandleFunc("", re.UpdatePermissionHandler).Methods(http.MethodPut).Name("updatepermission")
	// r.HandleFunc("/permission/{id}", permissionService.DeletePermissionHandler).Methods(http.MethodDelete).Name("deletepermission")
}
