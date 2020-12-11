package auth

import (
    "net/http"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func AuthMiddleware(hf http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        // to allow access from localhost
        w.Header().Set("Access-Control-Allow-Origin", "*")

        // AUTHORIZE IP
        if isIpAuthorized := isIpAuthorized(req); isIpAuthorized { // if ip is authorized, continue
            logger.InfoLog(" - AUTH -- IP is authorized.", req)

        } else { // if not authorized, send notification email and stop request
            logger.InfoLog(" - AUTH -- IP is unauthorized.", req)
            handleUnauthorizedIp(req)
            return
        }

        hf.ServeHTTP(w, req)
    })
}


