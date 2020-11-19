package btcxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetBalanceNodeXmlRpcRes struct {
    Response GetBalanceNodeXmlRpcResStruct
}
type GetBalanceNodeXmlRpcResStruct struct {
    Balance string
}

func (bs *BtcService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    res := model.GetBalanceRpcRes{ Balance: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := GetBalanceNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("getbalance", &rpcReq, &nodeRpcRes)

    if err == nil {
        res.Balance = nodeRpcRes.Response.Balance
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}