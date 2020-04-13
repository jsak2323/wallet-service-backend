package btc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func (bs *BtcService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    res := model.GetBalanceRpcRes{ Balance: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("getbalance", &rpcReq, &res)

    if err != nil { 
        return &res, err

    } else if res.Balance == "0" {
        return &res, errors.New("Unexpected error occured in Node.")

    } else {
        return &res, nil
    }
}