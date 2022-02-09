package xmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetBalanceXmlRpcRes struct {
	Content GetBalanceXmlRpcResStruct
}
type GetBalanceXmlRpcResStruct struct {
	Balance string
	Error   string
}

func (gs *GeneralService) GetBalance(ctx context.Context, rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
	res := model.GetBalanceRpcRes{Balance: "0"}

	rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
	client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcRes := GetBalanceXmlRpcRes{}

	err := client.XmlRpcCall(gs.Symbol+"Rpc.GetBalance", &rpcReq, &rpcRes)

	if err != nil {
		return &res, err

	} else if rpcRes.Content.Error != "" {
		return &res, errors.New(rpcRes.Content.Error)

	} else if rpcRes.Content.Balance == "0" || rpcRes.Content.Balance == "" {
		return &res, errors.New("Unexpected error occured in Node.")
	}

	res.Balance = rpcRes.Content.Balance
	return &res, nil
}
