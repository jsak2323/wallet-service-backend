package thetaxmlrpc

import (
    // "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    // "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type GetBalanceXmlRpcRes struct {
    Content GetBalanceXmlRpcResStruct
}
type GetBalanceXmlRpcResStruct struct {
    Balance string
    Error   string
}

func (ts *ThetaService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    res := model.GetBalanceRpcRes{ Balance: "0" }
    return &res, nil
}


