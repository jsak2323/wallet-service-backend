package main

import(
    "net/http"
    "database/sql"

    "github.com/gorilla/mux"
    
    h "github.com/btcid/wallet-services-backend/pkg/http/handlers"
    c "github.com/btcid/wallet-services-backend/pkg/http/handlers/cron"
    "github.com/btcid/wallet-services-backend/pkg/database/mysql"
)

func SetRoutes(r *mux.Router, mysqlDbConn *sql.DB) {

    // REPOSITORIES
    healthCheckRepo := mysql.NewMysqlHealthCheckRepository(mysqlDbConn)



    // REST ROUTES

    // getblockcount
    getBlockCountService := h.NewGetBlockCountService()
    r.HandleFunc("/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)



    // CRON ROUTES

    // healthcheck
    healthCheckService := c.NewHealthCheckService(healthCheckRepo)
    r.HandleFunc("/cron/healthcheck", healthCheckService.HealthCheckHandler).Methods(http.MethodGet)

}