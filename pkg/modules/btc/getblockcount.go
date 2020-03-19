package btc

import(
    // "fmt"
    "errors"    
    "strconv"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type BtcService struct {}

func (bs *BtcService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("getblockcount", &rpcReq, &res)

    return handleResponse(&res, err)
}

func (bs *BtcService) ConfirmBlockCount() (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    cryptoApisService := NewCryptoApisService()
    casRes, err := cryptoApisService.GetNodeInfo()

    if err == nil {
        res.Blocks = strconv.Itoa(casRes.Payload.Blocks)
    }

    return handleResponse(&res, err)
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