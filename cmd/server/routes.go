package main

import(
    "net/http"
    "database/sql"

    "github.com/gorilla/mux"
    
    h "github.com/btcid/wallet-services-backend/pkg/http/handlers"
    c "github.com/btcid/wallet-services-backend/pkg/http/handlers/cron"
    "github.com/btcid/wallet-services-backend/pkg/database/mysql"
    "github.com/btcid/wallet-services-backend/pkg/modules"
)

func SetRoutes(r *mux.Router, mysqlDbConn *sql.DB) {

    // REPOSITORIES
    healthCheckRepo := mysql.NewMysqlHealthCheckRepository(mysqlDbConn)


    // MODULE SERVICES
    ModuleServices := modules.NewModuleServices(healthCheckRepo)


    // REST ROUTES

    // getblockcount
    getBlockCountService := h.NewGetBlockCountService(ModuleServices)
    r.HandleFunc("/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)



    // CRON ROUTES

    // healthcheck
    healthCheckService := c.NewHealthCheckService(healthCheckRepo, ModuleServices)
    r.HandleFunc("/cron/healthcheck", healthCheckService.HealthCheckHandler).Methods(http.MethodGet)

}