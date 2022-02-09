package role

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type RoleHandler interface {
	CreateRoleHandler(w http.ResponseWriter, req *http.Request)
	CreateRolePermissionHandler(w http.ResponseWriter, req *http.Request)
	DeleteRoleHandler(w http.ResponseWriter, req *http.Request)
	DeleteRolePermissionHandler(w http.ResponseWriter, req *http.Request)
	ListRoleHandler(w http.ResponseWriter, req *http.Request)
	UpdateRoleHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewRoleHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
