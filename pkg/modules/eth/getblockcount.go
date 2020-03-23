package eth

import(
    "math"
    "errors"

    "github.com/ethereum/go-ethereum/common/hexutil"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

func (es *EthService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("EthRpc.GetBlockCount", &rpcReq, &res)

    if err != nil { 
        return &res, err

    } else if res.Blocks == "0" {
        return &res, errors.New("Unexpected error occured in Node.")

    } else {
        return &res, nil
    }
}

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