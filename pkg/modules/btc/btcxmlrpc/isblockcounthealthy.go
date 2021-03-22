package btcxmlrpc

import (
    "math"
    "strconv"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    modules_util "github.com/btcid/wallet-services-backend-go/pkg/modules/util"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/btc"
)

func (bs *BtcService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
    isBlockCountHealthy := false
    healthyBlockDiff    := config.CURR[bs.GetSymbol()].Config.HealthyBlockDiff
    blockDiff           := 0

    cryptoApisService := btc.NewCryptoApisService()
    getNodeInfoRes, err := cryptoApisService.GetNodeInfo()

    if err != nil { // if third party service fail, compare with previous blockcount
        logger.Log(" - "+bs.GetSymbol()+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" isBlockCountHealthy service err: "+err.Error()+". Executing fallback ...")
        isBlockCountFallbackHealthy, fallbackErr := modules_util.IsBlockCountHealthyFallback(bs, nodeBlockCount, rpcConfigId)
        if fallbackErr != nil { return isBlockCountFallbackHealthy, blockDiff, fallbackErr }
        isBlockCountHealthy = isBlockCountFallbackHealthy

    } else {
        blockDiff = nodeBlockCount - getNodeInfoRes.Payload.Blocks
        blockDiff = int(math.Abs(float64(blockDiff)))

        isBlockCountHealthy = blockDiff > healthyBlockDiff
    }

    return isBlockCountHealthy, blockDiff, nil
}
