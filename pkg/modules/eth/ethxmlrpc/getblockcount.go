package ethxmlrpc

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func (es *EthService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("EthRpc.GetBlockCount", &rpcReq, &res)

    if err == nil {
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}