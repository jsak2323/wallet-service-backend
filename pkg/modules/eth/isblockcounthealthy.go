package eth

import (
    "math"

    "github.com/ethereum/go-ethereum/common/hexutil"

    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

func (es *EthService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    healthyBlockDiff    := config.CURR["ETH"].Config.HealthyBlockDiff
    blockDiff           := 0

    infuraService := NewInfuraService()
    ethBlockNumberRes, err := infuraService.EthBlockNumber()

    if err != nil { // if third party service fail, compare with previous blockcount
        logger.Log(" - ETH isBlockCountHealthy  err: "+err.Error())
        previousHealthCheck, err := es.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
        if err != nil { return isBlockCountHealthy, blockDiff, err }

        if nodeBlockCount == previousHealthCheck.BlockCount { // if it's still the same, then it's not healthy
            isBlockCountHealthy = false
        }

    } else {
        ethBlockNumberHex := ethBlockNumberRes.Result
        ethBlockNumberUint64, _ := hexutil.DecodeUint64(ethBlockNumberHex)

        blockDiff = nodeBlockCount - int(ethBlockNumberUint64)
        blockDiff = int(math.Abs(float64(blockDiff)))

        isBlockCountHealthy = blockDiff <= healthyBlockDiff
    }

    return isBlockCountHealthy, blockDiff, nil
}