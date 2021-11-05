package xmlrpc

import (
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (gms *GeneralMapService) IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error) {
	isBlockCountHealthy := false
	PARENTSYMBOL := gms.ParentSymbol
	healthyBlockDiff := config.CURR[PARENTSYMBOL].Config.HealthyBlockDiff
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

	logger.Log(" - " + PARENTSYMBOL + " rpcConfigId: " + strconv.Itoa(rpcConfigId) + " nodeBlockCount: " + strconv.Itoa(nodeBlockCount) + ", previousBlockCount: " + strconv.Itoa(previousBlockCount))

	return isBlockCountHealthy, blockDiff, nil
}
