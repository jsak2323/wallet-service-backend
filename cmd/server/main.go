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
    mysqlDbConn := config.MysqlDbConn()
    defer mysqlDbConn.Close()
    
    r := mux.NewRouter()

    SetRoutes(r, mysqlDbConn)

    r.Use(mw.LogMiddleware)
    r.Use(mw.AuthMiddleware)

    server := &http.Server{
        Handler         : r,
        Addr            : ":"+config.CONF.Port,
        WriteTimeout    : 120 * time.Second,
        ReadTimeout     : 120 * time.Second,
    }

    fmt.Println()
    fmt.Println("Running server on localhost:"+config.CONF.Port)
    fmt.Println("\n\n\n")

    log.Fatal(server.ListenAndServe())
}