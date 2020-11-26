package xmlrpc

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetNewAddressXmlRpcRes struct {
    Content GetNewAddressXmlRpcResStruct
}
type GetNewAddressXmlRpcResStruct struct {
    Address string
    Error   string
}

func (gs *GeneralService) GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*model.GetNewAddressRpcRes, error) {
    res := model.GetNewAddressRpcRes{}

    rpcReq := util.GenerateRpcReq(rpcConfig, addressType, "", "")
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := GetNewAddressXmlRpcRes{}

    err := client.XmlRpcCall(gs.Symbol+"Rpc.GetNewAddress", &rpcReq, &rpcRes)

    if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" {
        return &res, errors.New(rpcRes.Content.Error)

    } else if rpcRes.Content.Address == "" {
        return &res, errors.New("Unexpected error occured in Node.")

    }

    res.Address = rpcRes.Content.Address
    return &res, nil
}

