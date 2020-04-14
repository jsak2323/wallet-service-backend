package btcxmlrpc

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func (es *EthService) ListTransactions(rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
    res := model.ListTransactionsRpcRes{ Transactions: "{\"transactions\": \"testeth\"}" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    return &res, nil

    // err := xmlrpc.XmlRpcCall("listtransactions", &rpcReq, &res)

    // if err != nil { 
    //     return &res, err

    // } else if res.Transactions == "" {
    //     return &res, errors.New("Unexpected error occured in Node.")

    // } else {
    //     return &res, nil
    // }
}