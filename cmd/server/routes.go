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

    // -- getblockcount
    getBlockCountService := h.NewGetBlockCountService(ModuleServices)
    r.HandleFunc("/blockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/blockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)

    // -- getbalance
    getBalanceService := h.NewGetBalanceService(ModuleServices)
    r.HandleFunc("/balance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/balance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)

    // -- listtransactions
    listTransactionsService := h.NewListTransactionsService(ModuleServices)
    r.HandleFunc("/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)

    // -- sendtoaddress
    sendToAddressService := h.NewSendToAddressService(ModuleServices)
    r.HandleFunc("/{symbol}/sendtoaddress/{address}/{amount}", sendToAddressService.SendToAddressHandler).Methods(http.MethodGet)


    // CRON ROUTES

    // -- healthcheck
    healthCheckService := c.NewHealthCheckService(healthCheckRepo, ModuleServices)
    r.HandleFunc("/cron/healthcheck", healthCheckService.HealthCheckHandler).Methods(http.MethodGet)

}