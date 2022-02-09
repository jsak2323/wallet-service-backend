package permission

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type PermissionHandler interface {
	ListPermissionHandler(w http.ResponseWriter, req *http.Request)
	CreatePermissionHandler(w http.ResponseWriter, req *http.Request)
	UpdatePermissionHandler(w http.ResponseWriter, req *http.Request)
	DeletePermissionHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewPermissionHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
