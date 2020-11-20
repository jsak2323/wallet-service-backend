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


    // API ROUTES

    // -- GET getblockcount
    getBlockCountService := h.NewGetBlockCountService(ModuleServices)
    r.HandleFunc("/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getblockcount", getBlockCountService.GetBlockCountHandler).Methods(http.MethodGet)

    // -- GET getbalance
    getBalanceService := h.NewGetBalanceService(ModuleServices)
    r.HandleFunc("/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getbalance", getBalanceService.GetBalanceHandler).Methods(http.MethodGet)

    // -- GET listtransactions
    listTransactionsService := h.NewListTransactionsService(ModuleServices)
    r.HandleFunc("/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/listtransactions", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/listtransactions/{limit}", listTransactionsService.ListTransactionsHandler).Methods(http.MethodGet)

    // -- GET getnewaddress
    getNewAddressService := h.NewGetNewAddressService(ModuleServices)
    r.HandleFunc("/{symbol}/getnewaddress", getNewAddressService.GetNewAddressHandler).Methods(http.MethodGet)
    r.HandleFunc("/{symbol}/getnewaddress/{type}", getNewAddressService.GetNewAddressHandler).Methods(http.MethodGet)

    // -- GET addresstype
    addressTypeService := h.NewAddressTypeService(ModuleServices)
    r.HandleFunc("/{symbol}/addresstype/{address}", addressTypeService.AddressTypeHandler).Methods(http.MethodGet)

    // -- POST sendtoaddress
    sendToAddressService := h.NewSendToAddressService(ModuleServices)
    r.HandleFunc("/sendtoaddress", sendToAddressService.SendToAddressHandler).Methods(http.MethodPost)
    /*
        curl example:
        curl --header "Content-Type: application/json" --request POST --data '{"symbol":"btc", "amount":"0.001", "address":"2MtU6EMx37AYrCNj1RcRr6bw66QqHYw4D4R"}' localhost:3000/sendtoaddress | jq
    */


    // CRON ROUTES

    // -- GET healthcheck
    healthCheckService := c.NewHealthCheckService(healthCheckRepo, ModuleServices)
    r.HandleFunc("/cron/healthcheck", healthCheckService.HealthCheckHandler).Methods(http.MethodGet)

}


