package auth

import (
    "strings"
    "net/http"

    "github.com/btcid/wallet-services-backend/cmd/config"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"    
)

func isIpAuthorized(req *http.Request) bool {
    ip := strings.Split(req.RemoteAddr, ":")[0]

    if isAuthorized, _ := util.InArray(ip, config.CONF.AuthorizedIps); !isAuthorized { 
        return false
    }
    return true
}


