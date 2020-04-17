package btcxmlrpc

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetBlockCountNodeXmlRpcRes struct {
    Blocks string
}

func (bs *BtcService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := GetBlockCountNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("getblockcount", &rpcReq, &nodeRpcRes)

    if err == nil {
        res.Blocks = nodeRpcRes.Blocks
        
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}