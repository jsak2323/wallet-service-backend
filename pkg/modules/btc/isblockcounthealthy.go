package btc

import (
    "math"

    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

func (bs *BtcService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    healthyBlockDiff    := config.CURR["BTC"].Config.HealthyBlockDiff
    blockDiff           := 0

    cryptoApisService := NewCryptoApisService()
    getNodeInfoRes, err := cryptoApisService.GetNodeInfo()

    if err != nil { // if third party service fail, compare with previous blockcount
        logger.Log(" - BTC isBlockCountHealthy cryptoApisService.GetNodeInfo() err: "+err.Error())
        previousHealthCheck, err := bs.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
        if err != nil { return isBlockCountHealthy, blockDiff, err }

        if nodeBlockCount == previousHealthCheck.BlockCount { // if it's still the same, then it's not healthy
            isBlockCountHealthy = false
        }

    } else {
        blockDiff = nodeBlockCount - getNodeInfoRes.Payload.Blocks
        blockDiff = int(math.Abs(float64(blockDiff)))

        isBlockCountHealthy = blockDiff <= healthyBlockDiff
    }

    return isBlockCountHealthy, blockDiff, nil
}
