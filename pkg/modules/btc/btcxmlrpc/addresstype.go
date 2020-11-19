package btcxmlrpc

import (
    // "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    // "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

// todo: addresstype function for btc (btc-middleware doesn't have the function yet)
func (bs *BtcService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{ AddressType: "" }

    return &res, nil
    
    // rpcReq := util.GenerateRpcReq(rpcConfig, address, "", "")
    // xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)


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