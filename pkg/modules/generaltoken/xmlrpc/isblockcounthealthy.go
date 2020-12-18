package xmlrpc

import (
    "strconv"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

func (gts *GeneralTokenService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    PARENTSYMBOL        := gts.ParentSymbol
    healthyBlockDiff    := config.CURR[PARENTSYMBOL].Config.HealthyBlockDiff
    blockDiff           := 0

    previousHealthCheck, err := gts.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
    if err != nil { return isBlockCountHealthy, blockDiff, err }

    previousBlockCount := previousHealthCheck.BlockCount

    blockDiff = nodeBlockCount - previousBlockCount
    if blockDiff < healthyBlockDiff { 
        isBlockCountHealthy = true
    }

    logger.Log(" - "+PARENTSYMBOL+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" nodeBlockCount: "+strconv.Itoa(nodeBlockCount)+", previousBlockCount: "+strconv.Itoa(previousBlockCount))

    return isBlockCountHealthy, blockDiff, nil
}

