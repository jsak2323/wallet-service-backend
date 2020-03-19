package eth

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type EthService struct {}

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

func (es *EthService) ConfirmBlockCount() (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "100" }

    return &res, nil
}