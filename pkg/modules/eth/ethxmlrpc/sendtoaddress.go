package ethxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (es *EthService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
    res := model.SendToAddressRpcRes{}

    rpcReq := util.GenerateRpcReq(rpcConfig, address, amountInDecimal, "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("EthRpc.SendTransaction", &rpcReq, &res)

    if err == nil {
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}