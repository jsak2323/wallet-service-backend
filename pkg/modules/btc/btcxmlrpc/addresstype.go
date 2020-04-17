package btcxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

// todo: addresstype function for btc (btc-middleware doesn't have the function yet)
func (bs *BtcService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{ AddressType: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, address, "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    return &res, nil

    // nodeRpcRes := GetBalanceNodeXmlRpcRes{}

    // err := xmlrpc.XmlRpcCall("getbalance", &rpcReq, &nodeRpcRes)

    // if err == nil {
    //     res.Balance = nodeRpcRes.Response.Balance
    //     return &res, nil

    // } else if err != nil {
    //     return &res, err

    // } else {
    //     return &res, errors.New("Unexpected error occured in Node.")
    // }
}