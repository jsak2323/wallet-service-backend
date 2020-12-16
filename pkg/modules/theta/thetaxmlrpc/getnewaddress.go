package thetaxmlrpc

import (
    // "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    // "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type GetNewAddressXmlRpcRes struct {
    Content GetNewAddressXmlRpcResStruct
}
type GetNewAddressXmlRpcResStruct struct {
    Address string
    Error   string
}

func (ts *ThetaService) GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*model.GetNewAddressRpcRes, error) {
    res := model.GetNewAddressRpcRes{}
    return &res, nil
}


