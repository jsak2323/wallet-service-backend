package handlers

import(
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func GetBlockCountHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Symbol: %v\n", vars["symbol"])
}