package xmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetBalanceXmlRpcRes struct {
    Content GetBalanceXmlRpcResStruct
}
type GetBalanceXmlRpcResStruct struct {
    Balance string
    Error   string
}

func (gs *GeneralService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    res := model.GetBalanceRpcRes{ Balance: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := GetBalanceXmlRpcRes{}

    err := client.XmlRpcCall(gs.Symbol+"Rpc.GetBalance", &rpcReq, &rpcRes)

    if err == nil {
        res.Balance = rpcRes.Content.Balance
        return &res, nil

    } else if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" {
        return &res, errors.New(rpcRes.Content.Error)

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}


