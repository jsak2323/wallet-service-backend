package thetaxmlrpc

import (
    "strconv"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

func (ts *ThetaService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    PARENTSYMBOL        := ts.GetParentSymbol()
    healthyBlockDiff    := config.CURR[PARENTSYMBOL].Config.HealthyBlockDiff
    blockDiff           := 0

    previousHealthCheck, err := ts.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
    if err != nil { return isBlockCountHealthy, blockDiff, err }

    previousBlockCount := previousHealthCheck.BlockCount

    blockDiff = nodeBlockCount - previousBlockCount
    if blockDiff < healthyBlockDiff { 
        isBlockCountHealthy = true
    }

    logger.Log(" - "+PARENTSYMBOL+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" nodeBlockCount: "+strconv.Itoa(nodeBlockCount)+", previousBlockCount: "+strconv.Itoa(previousBlockCount))

    return isBlockCountHealthy, blockDiff, nil
}


