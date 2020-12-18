package xmlrpc

import (
    "errors"
    "strings"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type GetBalanceXmlRpcRes struct {
    Content GetBalanceXmlRpcResStruct
}
type GetBalanceXmlRpcResStruct struct {
    Balance string
    Error   string
}

func (gts *GeneralTokenService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    res := model.GetBalanceRpcRes{ Balance: "0" }

    token := strings.ToLower(gts.Symbol)

    rpcReq := util.GenerateRpcReq(rpcConfig, token, "", "")
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := GetBalanceXmlRpcRes{}

    err := client.XmlRpcCall(gs.ParentSymbol+"Rpc.GetBalance", &rpcReq, &rpcRes)

    if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" {
        return &res, errors.New(rpcRes.Content.Error)

    } else if rpcRes.Content.Balance == "0" || rpcRes.Content.Balance == "" {
        return &res, errors.New("Unexpected error occured in Node.")
    } 

    res.Balance = rpcRes.Content.Balance
    return &res, nil
}


