package xmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type ListTransactionsXmlRpcRes struct {
    Content ListTransactionsXmlRpcResStruct
}
type ListTransactionsXmlRpcResStruct struct {
    Transactions string
    Error        string
}

func (gs *GeneralService) ListTransactions(rpcConfig rc.RpcConfig) (*model.ListTransactionsRpcRes, error) {
    res := model.ListTransactionsRpcRes{}

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := ListTransactionsXmlRpcRes{}

    err := client.XmlRpcCall(gs.Symbol+"Rpc.ListTransactions", &rpcReq, &rpcRes)

    if err == nil {
        res.Transactions = rpcRes.Content.Transactions
        return &res, nil

    } else if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" {
        return &res, errors.New(rpcRes.Content.Error)

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}


