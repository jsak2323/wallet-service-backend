package main

import(
    "net/http"

    "github.com/gorilla/mux"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Key: %v\n", vars["key"])
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/products/{key}", ProductHandler)
}