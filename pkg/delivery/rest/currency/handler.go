package currency

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type CurrencyHandler interface {
	ListHandler(w http.ResponseWriter, req *http.Request)
	CreateHandler(w http.ResponseWriter, req *http.Request)
	UpdateHandler(w http.ResponseWriter, req *http.Request)
	CreateRpcHandler(w http.ResponseWriter, req *http.Request)
	DeleteRpcHandler(w http.ResponseWriter, req *http.Request)
	DeactivateHandler(w http.ResponseWriter, req *http.Request)
	ActivateHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewCurrencyHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
