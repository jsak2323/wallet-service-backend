package util

import (
    "fmt"
    "strconv"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func IsBlockCountHealthyFallback(service modules_m.ModuleService, nodeBlockCount int, rpcConfigId int) (bool, error) {
    isBlockCountHealthy := false 

    previousHealthCheck, err := service.GetHealthCheckRepo().GetByRpcConfigId(rpcConfigId)
    if err != nil { return isBlockCountHealthy, err }

    previousBlockCount := previousHealthCheck.BlockCount
    fmt.Println("nodeBlockCount: "+strconv.Itoa(nodeBlockCount))
    fmt.Println("previousBlockCount: "+strconv.Itoa(previousBlockCount))
    if nodeBlockCount > previousBlockCount { // if it's moved up from previous blockcount, then it's healthy
        isBlockCountHealthy = true
    }

    logger.Log(" - "+service.GetSymbol()+
        " rpcConfigId: "+strconv.Itoa(rpcConfigId)+
        " isBlockCountHealthy fallback -"+
        " previousBlockCount: "+strconv.Itoa(previousBlockCount)+
        ", nodeBlockCount: "+strconv.Itoa(nodeBlockCount)+
        ", isBlockCountHealthy: "+strconv.FormatBool(isBlockCountHealthy))

    return isBlockCountHealthy, nil
}