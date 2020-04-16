package ethxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func (es *EthService) SendToAddress(rpcConfig rc.RpcConfig, address string, amountInDecimal string) (*model.SendToAddressRpcRes, error) {
    res := model.SendToAddressRpcRes{ TxHash: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, address, amountInDecimal, "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("EthRpc.SendTransaction", &rpcReq, &res)

    if err == nil {
        res.TxHash = res.TxHash
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}