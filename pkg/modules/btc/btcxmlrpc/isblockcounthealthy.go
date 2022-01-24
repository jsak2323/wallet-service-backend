package btcxmlrpc

import (
	"context"
	"math"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/btc"
	modules_util "github.com/btcid/wallet-services-backend-go/pkg/modules/util"
)

func (bs *BtcService) IsBlockCountHealthy(ctx context.Context, nodeBlockCount int, rpcConfigId int) (bool, int, error) {
	isBlockCountHealthy := false
	healthyBlockDiff := config.CURRRPC[bs.currencyConfigId].Config.HealthyBlockDiff
	blockDiff := 0

	cryptoApisService := btc.NewCryptoApisService()
	getNodeInfoRes, err := cryptoApisService.GetNodeInfo()

	if err != nil { // if third party service fail, compare with previous blockcount
		logger.Log(" - "+bs.GetSymbol()+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" isBlockCountHealthy service err: "+err.Error()+". Executing fallback ...", ctx)
		isBlockCountFallbackHealthy, fallbackErr := modules_util.IsBlockCountHealthyFallback(ctx, bs, nodeBlockCount, rpcConfigId)
		if fallbackErr != nil {
			return isBlockCountFallbackHealthy, blockDiff, fallbackErr
		}
		isBlockCountHealthy = isBlockCountFallbackHealthy

	} else {
		blockDiff = nodeBlockCount - getNodeInfoRes.Payload.Blocks
		blockDiff = int(math.Abs(float64(blockDiff)))

		isBlockCountHealthy = blockDiff > healthyBlockDiff
	}

	return isBlockCountHealthy, blockDiff, nil
}
