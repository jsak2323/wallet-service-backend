package ethxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func (es *EthService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{}

    rpcReq := util.GenerateRpcReq(rpcConfig, address, "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("EthRpc.AddressType", &rpcReq, &res)

    if err == nil {
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}