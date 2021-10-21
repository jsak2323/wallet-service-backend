package xmlrpc

import (
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (gts *GeneralTokenMapService) AddressType(rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
	res := model.AddressTypeRpcRes{AddressType: ""}

	return &res, nil
}
