package logging

import(
    "fmt"
    "net/http"

    logger "github.com/btcid/wallet-services-backend/pkg/logging"
)

func LogMiddleware(hf http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.InfoLog(r.URL.String() + " hit. ", r)
        hf.ServeHTTP(w, r)
        logger.InfoLog(r.URL.String() + " done. ", r)
        fmt.Println()
    })
}


