package btcxmlrpc

import (
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (bs *BtcService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{ AddressType: "" }

    return &res, nil
}


