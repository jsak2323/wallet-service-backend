package xmlrpc

import (
	"context"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (gms *GeneralMapService) IsBlockCountHealthy(ctx context.Context, nodeBlockCount int, rpcConfigId int) (bool, int, error) {
	isBlockCountHealthy := false

	parentCurrency, err := config.GetCurrencyBySymbolTokenType(gms.ParentSymbol, cc.MainTokenType)
	if err != nil {
		return false, 0, err
	}

	healthyBlockDiff := parentCurrency.HealthyBlockDiff
	blockDiff := 0

	previousHealthCheck, err := gms.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
	if err != nil {
		return isBlockCountHealthy, blockDiff, err
	}

	previousBlockCount := previousHealthCheck.BlockCount

	blockDiff = nodeBlockCount - previousBlockCount
	if blockDiff > healthyBlockDiff {
		isBlockCountHealthy = true
	}

	logger.Log(" - "+gms.ParentSymbol+" rpcConfigId: "+strconv.Itoa(rpcConfigId)+" nodeBlockCount: "+strconv.Itoa(nodeBlockCount)+", previousBlockCount: "+strconv.Itoa(previousBlockCount), ctx)

	return isBlockCountHealthy, blockDiff, nil
}
