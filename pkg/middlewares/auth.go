package middlewares

import (
    "fmt"
    "net/http"
)

func AuthMiddleware(hf http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // todo: authorize requests
        fmt.Println("--- auth middleware")
        hf.ServeHTTP(w, r)
    })    
}