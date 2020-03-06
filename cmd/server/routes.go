package main

import(
    "github.com/gorilla/mux"

    "github.com/btcid/wallet-services-backend/pkg/handlers"
)

func SetRoutes(r *mux.Router) {
    r.HandleFunc("/getblockcount", handlers.GetBlockCountHandler)
    r.HandleFunc("/getblockcount/{symbol}", handlers.GetBlockCountHandler)
}