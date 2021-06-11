package main

import(
    "fmt"
    "log"
    "time"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"

    "github.com/go-redis/redis/v8"
    logm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/logging"
	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

func main() {
    mysqlDbConn := config.MysqlDbConn()
    defer mysqlDbConn.Close()

    redis := redis.NewClient(&redis.Options{
        Addr:     config.CONF.RedisHost,
    })
    
    r := mux.NewRouter()

	SetRoutes(r, mysqlDbConn, redis)

    corsOpts := cors.New(cors.Options{
        AllowedMethods: []string{
            http.MethodGet,
            http.MethodPost,
            http.MethodPut,
            http.MethodDelete,
            http.MethodOptions,
            http.MethodHead,
        },
        AllowedHeaders: []string{
            "*",
        },
    })

	r.Use(logm.LogMiddleware)

    server := &http.Server{
        Handler         : corsOpts.Handler(r),
        Addr            : ":"+config.CONF.Port,
        WriteTimeout    : 120 * time.Second,
        ReadTimeout     : 120 * time.Second,
    }

    fmt.Println()
    fmt.Println("Running server on localhost:"+config.CONF.Port)
    fmt.Println("\n\n\n")

    log.Fatal(server.ListenAndServe())
}
