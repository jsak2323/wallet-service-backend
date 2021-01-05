package xmlrpc

import (
    "fmt"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

func (gts *GeneralTokenService) GetLog(rpcConfig rc.RpcConfig, date string) (string, error) {
    
    fmt.Println("getLog hit")

}


