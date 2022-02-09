package ethxmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (es *EthService) GetNewAddress(ctx context.Context, rpcConfig rc.RpcConfig, addressType string) (*model.GetNewAddressRpcRes, error) {
	res := model.GetNewAddressRpcRes{}

	rpcReq := util.GenerateRpcReq(rpcConfig, addressType, "", "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	err := xmlrpc.XmlRpcCall("EthRpc.NewAccount", &rpcReq, &res)

	if err == nil {
		return &res, nil

	} else if err != nil {
		return &res, err

	} else {
		return &res, errors.New("Unexpected error occured in Node.")
	}
}
