package xmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type SendToAddressXmlRpcRes struct {
	Content SendToAddressXmlRpcResStruct
}
type SendToAddressXmlRpcResStruct struct {
	TxHash string
	Error  string
}

func (gs *GeneralService) SendToAddress(ctx context.Context, rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
	res := model.SendToAddressRpcRes{}

	rpcReq := util.GenerateRpcReq(rpcConfig, amountInDecimal, address, memo)
	client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcRes := SendToAddressXmlRpcRes{}

	err := client.XmlRpcCall(gs.Symbol+"Rpc.SendToAddress", &rpcReq, &rpcRes)

	if err != nil {
		return &res, err

	} else if rpcRes.Content.Error != "" {
		return &res, errors.New(rpcRes.Content.Error)

	} else if rpcRes.Content.TxHash == "" {
		return &res, errors.New("Unexpected error occured in Node.")
	}

	res.TxHash = rpcRes.Content.TxHash
	return &res, nil
}
