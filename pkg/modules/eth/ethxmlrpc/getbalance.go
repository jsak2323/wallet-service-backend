package ethxmlrpc

import (
	"context"
	"errors"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (es *EthService) GetBalance(ctx context.Context, rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
	balanceRes := struct{ Value string }{}
	res := model.GetBalanceRpcRes{Balance: "0"}

	rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	err := xmlrpc.XmlRpcCall("balance", &rpcReq, &balanceRes)

	res.Balance = balanceRes.Value

	if err == nil {
		return &res, nil

	} else if err != nil {
		return &res, err

	} else {
		return &res, errors.New("Unexpected error occured in Node.")
	}
}
