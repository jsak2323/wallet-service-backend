package btc

import(
    "fmt"
    "math"
    "errors"    
    // "strconv"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

func (bs *BtcService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("getblockcount", &rpcReq, &res)

    return handleResponse(&res, err)
}

func (bs *BtcService) IsBlockCountHealthy(nodeBlockCount int) (bool, int, error) {
    isBlockCountHealthy := false
    healthyBlockDiff    := config.CURR["BTC"].Config.HealthyBlockDiff
    blockDiff           := 0

    cryptoApisService := NewCryptoApisService()
    casRes, err := cryptoApisService.GetNodeInfo()

    if err != nil { // if third party service fail, compare with previous blockcount
        fmt.Println("err: "+err.Error())

    } else {
        blockDiff = nodeBlockCount - casRes.Payload.Blocks
        blockDiff = int(math.Abs(float64(blockDiff)))

        isBlockCountHealthy = blockDiff <= healthyBlockDiff
    }

    return isBlockCountHealthy, blockDiff, nil
}

func handleResponse(res *model.GetBlockCountRpcRes, err error) (*model.GetBlockCountRpcRes, error) {
    if err != nil { 
        return res, err

    } else if res.Blocks == "0" {
        return res, errors.New("Unexpected error occured in Node.")

    } else {
        return res, nil
    }
}