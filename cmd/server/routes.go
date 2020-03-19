package main

import(
    "net/http"

    "github.com/gorilla/mux"
    
    "github.com/btcid/wallet-services-backend/pkg/http/handlers"
    "github.com/btcid/wallet-services-backend/pkg/http/handlers/cron"
)

func SetRoutes(r *mux.Router) {
    r.HandleFunc("/getblockcount", handlers.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getblockcount", handlers.GetBlockCountHandler).Methods(http.MethodGet)
}

func SetCronRoutes(r *mux.Router) {
    r.HandleFunc("/cron/healthcheck", cron.HealthCheckHandler).Methods(http.MethodGet)
}

