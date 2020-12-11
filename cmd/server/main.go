package main

import(
    "fmt"
    "log"
    "time"
    "net/http"

    "github.com/gorilla/mux"

    logm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/logging"
    authm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/auth"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

func main() {
    mysqlDbConn := config.MysqlDbConn()
    defer mysqlDbConn.Close()
    
    r := mux.NewRouter()

    SetRoutes(r, mysqlDbConn)

    r.Use(logm.LogMiddleware)
    r.Use(authm.AuthMiddleware)

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