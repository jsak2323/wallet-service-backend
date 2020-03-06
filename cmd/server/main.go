package main

import(
    "fmt"
    "log"
    "time"
    "net/http"

    "github.com/gorilla/mux"

    "github.com/btcid/wallet-services-backend/cmd/config"
)

func main() {
    r := mux.NewRouter()
    SetRoutes(r)

    server := &http.Server{
        Handler         : r,
        Addr            : "127.0.0.1:"+config.CONF.Port,
        WriteTimeout    : 15 * time.Second,
        ReadTimeout     : 15 * time.Second,
    }

    fmt.Println("Running server on localhost:"+config.CONF.Port)

    log.Fatal(server.ListenAndServe())
}