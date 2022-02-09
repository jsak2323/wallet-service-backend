package ethxmlrpc

import (
	// "errors"
	"context"
	"strconv"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (es *EthService) ListTransactions(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
	res := model.ListTransactionsRpcRes{}

	rpcReq := util.GenerateRpcReq(rpcConfig, strconv.Itoa(limit), "", "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	err := xmlrpc.XmlRpcCall("EthRpc.ListTransactions", &rpcReq, &res)

	if err != nil {
		return &res, err

	} else {
		return &res, nil
	}
	// else if res.Transactions == "" {
	//     return &res, errors.New("Unexpected error occured in Node.")

	// }

}
