package xmlrpc

import (
    "strconv"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

func (gs *GeneralService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    SYMBOL              := gs.GetSymbol()
    healthyBlockDiff    := config.CURR[SYMBOL].Config.HealthyBlockDiff
    blockDiff           := 0

    logMsg := " - "+SYMBOL+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" nodeBlockCount: "+strconv.Itoa(nodeBlockCount)+", previousBlockCount: "+strconv.Itoa(previousBlockCount)

    if healthyBlockDiff == 0 {
        logger.Log(logMsg)
        return true, blockDiff, nil
    }

    previousHealthCheck, err := gs.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
    if err != nil { return isBlockCountHealthy, blockDiff, err }

    previousBlockCount := previousHealthCheck.BlockCount

    blockDiff = nodeBlockCount - previousBlockCount
    if blockDiff > healthyBlockDiff { 
        isBlockCountHealthy = true
    }

    logger.Log(logMsg)
    return isBlockCountHealthy, blockDiff, nil
}


