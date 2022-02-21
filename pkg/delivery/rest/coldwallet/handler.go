package coldwallet

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type ColdWalletHandler interface {
	ActivateHandler(w http.ResponseWriter, req *http.Request)
	DeactivateHandler(w http.ResponseWriter, req *http.Request)
	ListHandler(w http.ResponseWriter, req *http.Request)
	CreateHandler(w http.ResponseWriter, req *http.Request)
	UpdateHandler(w http.ResponseWriter, req *http.Request)
	UpdateBalanceHandler(w http.ResponseWriter, req *http.Request)
	SendToHotHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewColdWalletHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
