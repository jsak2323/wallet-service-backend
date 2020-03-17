package main

import(
    "net/http"

    "github.com/gorilla/mux"
    
    "github.com/btcid/wallet-services-backend/pkg/handlers"
)

func SetRoutes(r *mux.Router) {
    r.HandleFunc("/getblockcount", handlers.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getblockcount", handlers.GetBlockCountHandler).Methods(http.MethodGet)
}