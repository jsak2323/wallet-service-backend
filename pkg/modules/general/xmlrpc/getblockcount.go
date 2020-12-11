package xmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type GetBlockCountXmlRpcRes struct {
    Content GetBlockCountXmlRpcResStruct
}
type GetBlockCountXmlRpcResStruct struct {
    Blocks  string
    Error   string
}

func (gs *GeneralService) GetBlockCount(rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
    res := model.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := GetBlockCountXmlRpcRes{}

    err := client.XmlRpcCall(gs.Symbol+"Rpc.GetBlockCount", &rpcReq, &rpcRes)

    if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" {
        return &res, errors.New(rpcRes.Content.Error)

    } else if rpcRes.Content.Blocks == "0" || rpcRes.Content.Blocks == "" {
        return &res, errors.New("Unexpected error occured in Node.")
    }
    
    res.Blocks = rpcRes.Content.Blocks
    return &res, nil
}


