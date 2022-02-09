package btcxmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type SendToAddressNodeXmlRpcRes struct {
	Content SendToAddressNodeXmlRpcResStruct
}
type SendToAddressNodeXmlRpcResStruct struct {
	Tx string
}

func (bs *BtcService) SendToAddress(ctx context.Context, rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
	res := model.SendToAddressRpcRes{TxHash: ""}

	rpcReq := util.GenerateRpcReq(rpcConfig, address, amountInDecimal, "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	nodeRpcRes := SendToAddressNodeXmlRpcRes{}

	err := xmlrpc.XmlRpcCall("sendtoaddress", &rpcReq, &nodeRpcRes)

	if err == nil {
		res.TxHash = nodeRpcRes.Content.Tx
		return &res, nil

	} else if err != nil {
		return &res, err

	} else {
		return &res, errors.New("Unexpected error occured in Node.")
	}
}
