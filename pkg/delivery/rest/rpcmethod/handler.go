package rpcmethod

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type RpcMethodHandler interface {
	ListHandler(w http.ResponseWriter, req *http.Request)
	GetByRpcConfigIdHandler(w http.ResponseWriter, req *http.Request)
	CreateHandler(w http.ResponseWriter, req *http.Request)
	UpdateHandler(w http.ResponseWriter, req *http.Request)
	DeleteHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewRpcMethodHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
