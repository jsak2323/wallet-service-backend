package middlewares

import (
    "net/http"
)

func AuthMiddleware(hf http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // todo: authorize requests

        // to allow access from localhost
        w.Header().Set("Access-Control-Allow-Origin", "*")

        hf.ServeHTTP(w, r)
    })
}


