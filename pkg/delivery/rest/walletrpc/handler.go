package walletrpc

import (
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/service"
	"github.com/gorilla/mux"
)

type WalletRpcHandler interface {
	AddressTypeHandler(w http.ResponseWriter, req *http.Request)
	GetBalanceHandler(w http.ResponseWriter, req *http.Request)
	GetBlockCountHandler(w http.ResponseWriter, req *http.Request)
	GetHealthCheckHandler(w http.ResponseWriter, req *http.Request)
	GetLogHandler(w http.ResponseWriter, req *http.Request)
	GetNewAddressHandler(w http.ResponseWriter, req *http.Request)
	ListTransactionsHandler(w http.ResponseWriter, req *http.Request)
	ListWithdrawsHandler(w http.ResponseWriter, req *http.Request)
	SendToAddressHandler(w http.ResponseWriter, req *http.Request)
	MaintenanceListHandler(w http.ResponseWriter, req *http.Request)
}

type (
	Rest struct {
		routes *mux.Router
		svc    service.Service
	}
)

func NewWalletRpcHandler(
	routes *mux.Router,
	svc service.Service,
) *Rest {
	return &Rest{
		routes: routes,
		svc:    svc,
	}
}
