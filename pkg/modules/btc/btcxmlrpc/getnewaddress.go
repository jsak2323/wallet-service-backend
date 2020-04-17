package btcxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetNewAddressNodeXmlRpcRes struct {
    Response GetNewAddressNodeXmlRpcResStruct
}
type GetNewAddressNodeXmlRpcResStruct struct {
    Address string
}

func (bs *BtcService) GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*model.GetNewAddressRpcRes, error) {
    res := model.GetNewAddressRpcRes{ Address: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, addressType, "", "")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := GetNewAddressNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("getnewaddress", &rpcReq, &nodeRpcRes)

    if err == nil {
        res.Address = nodeRpcRes.Response.Address
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}