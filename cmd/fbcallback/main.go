package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	logm "github.com/btcid/wallet-services-backend-go/pkg/middlewares/logging"
)

func main() {
	mysqlDbConn := config.MysqlDbConn()
	defer mysqlDbConn.Close()

	r := mux.NewRouter()

	fireblocksService := h.NewFireblocksService()
	r.HandleFunc("/fireblocks/tx_sign_request", fireblocksService.CallbackHandler).Methods(http.MethodPost).Name("fireblockscallback")

	corsOpts := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})

	r.Use(logm.LogMiddleware)

	server := &http.Server{
		Handler:      corsOpts.Handler(r),
		Addr:         ":" + config.CONF.FireblocksCallbackPort,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	fmt.Println()
	fmt.Println("Running server on localhost:"+config.CONF.FireblocksCallbackPort)
	fmt.Println("\n\n\n")
	
	log.Fatal(server.ListenAndServeTLS(config.CONF.FireblocksCallbackSSLCert, config.CONF.FireblocksCallbackSSLKey))
}
