package btcxmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetBlockCountNodeXmlRpcRes struct {
	Blocks string
}

func (bs *BtcService) GetBlockCount(ctx context.Context, rpcConfig rc.RpcConfig) (*model.GetBlockCountRpcRes, error) {
	res := model.GetBlockCountRpcRes{Blocks: "0"}

	rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	nodeRpcRes := GetBlockCountNodeXmlRpcRes{}

	err := xmlrpc.XmlRpcCall("getblockcount", &rpcReq, &nodeRpcRes)

	if err == nil {
		res.Blocks = nodeRpcRes.Blocks

		return &res, nil

	} else if err != nil {
		return &res, err

	} else {
		return &res, errors.New("Unexpected error occured in Node.")
	}
}
