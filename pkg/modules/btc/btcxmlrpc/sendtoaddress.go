package btcxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type SendToAddressNodeXmlRpcRes struct {
    Response SendToAddressNodeXmlRpcResStruct
}
type SendToAddressNodeXmlRpcResStruct struct {
    TxHash string
}

func (bs *BtcService) SendToAddress(rpcConfig rc.RpcConfig, address string, amountInDecimal string) (*model.SendToAddressRpcRes, error) {
    res := model.SendToAddressRpcRes{ TxHash: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, address, amountInDecimal, "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := SendToAddressNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("SendToAddress", &rpcReq, &nodeRpcRes)

    if err == nil {
        res.Balance = nodeRpcRes.Response.Balance
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}