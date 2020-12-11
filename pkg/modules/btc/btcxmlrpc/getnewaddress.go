package btcxmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type GetNewAddressNodeXmlRpcRes struct {
    Content GetNewAddressNodeXmlRpcResStruct
}
type GetNewAddressNodeXmlRpcResStruct struct {
    Address string
}

func (bs *BtcService) GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*model.GetNewAddressRpcRes, error) {
    res := model.GetNewAddressRpcRes{ Address: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, addressType, "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := GetNewAddressNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("getnewaddress", &rpcReq, &nodeRpcRes)

    if err == nil {
        res.Address = nodeRpcRes.Content.Address
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}