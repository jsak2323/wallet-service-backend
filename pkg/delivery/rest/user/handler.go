package user

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type UserHandler interface {
	LoginHandler(w http.ResponseWriter, req *http.Request)
	LogoutHandler(w http.ResponseWriter, req *http.Request)
	ActivateUserHandler(w http.ResponseWriter, req *http.Request)
	DeactivateUserHandler(w http.ResponseWriter, req *http.Request)
	ListUserHandler(w http.ResponseWriter, req *http.Request)
	CreateUserHandler(w http.ResponseWriter, req *http.Request)
	UpdateUserHandler(w http.ResponseWriter, req *http.Request)
	AddUserRolesHandler(w http.ResponseWriter, req *http.Request)
	DeleteUserRoleHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewUserHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
