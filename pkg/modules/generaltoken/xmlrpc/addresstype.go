package xmlrpc

import (
	"context"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (gts *GeneralTokenService) AddressType(ctx context.Context, rpcConfig rc.RpcConfig, address string) (*model.AddressTypeRpcRes, error) {
	res := model.AddressTypeRpcRes{AddressType: ""}

	return &res, nil
}
