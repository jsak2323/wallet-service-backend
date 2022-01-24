package ethxmlrpc

import (
	"context"
	"math"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/eth"
	modules_util "github.com/btcid/wallet-services-backend-go/pkg/modules/util"
)

func (es *EthService) IsBlockCountHealthy(ctx context.Context, nodeBlockCount int, rpcConfigId int) (bool, int, error) {
	isBlockCountHealthy := false
	healthyBlockDiff := config.CURRRPC[es.currencyConfigId].Config.HealthyBlockDiff
	blockDiff := 0

	infuraService := eth.NewInfuraService()
	ethBlockNumberRes, err := infuraService.EthBlockNumber()

	if err != nil { // if third party service fail, compare with previous blockcount
		logger.Log(" - "+es.GetSymbol()+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" isBlockCountHealthy service err: "+err.Error()+". Executing fallback ...", ctx)
		isBlockCountFallbackHealthy, fallbackErr := modules_util.IsBlockCountHealthyFallback(ctx, es, nodeBlockCount, rpcConfigId)
		if fallbackErr != nil {
			return isBlockCountFallbackHealthy, blockDiff, fallbackErr
		}
		isBlockCountHealthy = isBlockCountFallbackHealthy

	} else {
		ethBlockNumberHex := ethBlockNumberRes.Result
		ethBlockNumberUint64, _ := hexutil.DecodeUint64(ethBlockNumberHex)

		blockDiff = nodeBlockCount - int(ethBlockNumberUint64)
		blockDiff = int(math.Abs(float64(blockDiff)))

		isBlockCountHealthy = blockDiff > healthyBlockDiff
	}

	return isBlockCountHealthy, blockDiff, nil
}
