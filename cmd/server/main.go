package main

import(
    "fmt"
    "log"
    "time"
    "net/http"

    "github.com/gorilla/mux"

    mw "github.com/btcid/wallet-services-backend/pkg/middlewares"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

func main() {
    r := mux.NewRouter()
    SetRoutes(r)
    
    r.Use(mw.LogMiddleware)
    r.Use(mw.AuthMiddleware)

    server := &http.Server{
        Handler         : r,
        Addr            : "127.0.0.1:"+config.CONF.Port,
        WriteTimeout    : 15 * time.Second,
        ReadTimeout     : 15 * time.Second,
    }

    fmt.Println()
    fmt.Println("Running server on localhost:"+config.CONF.Port)
    fmt.Println("\n\n\n")

    log.Fatal(server.ListenAndServe())
}