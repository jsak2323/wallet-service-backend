package thetaxmlrpc

import (
    // "errors"
    // "strconv"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    // "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type ListTransactionsXmlRpcRes struct {
    Content ListTransactionsXmlRpcResStruct
}
type ListTransactionsXmlRpcResStruct struct {
    Transactions string
    Error        string
}

// to be implemented
func (ts *ThetaService) ListTransactions(rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
    res := model.ListTransactionsRpcRes{}
    return &res, nil
}


