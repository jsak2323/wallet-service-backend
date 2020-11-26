package xmlrpc

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
)

func (gs *GeneralService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{ AddressType: "" }

    return &res, nil
}


