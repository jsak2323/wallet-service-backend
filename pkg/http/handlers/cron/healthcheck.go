package cron

import(
    "fmt"
    "net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) { 
    fmt.Println("HealthCheckHandler hit.")
}
