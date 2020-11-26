package btcxmlrpc

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
)

// todo: addresstype function for btc (btc-middleware doesn't have the function yet)
func (bs *BtcService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
    res := model.AddressTypeRpcRes{ AddressType: "" }

    return &res, nil
}