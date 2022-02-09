package xmlrpc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type ListTransactionsXmlRpcRes struct {
	Content ListTransactionsXmlRpcResStruct
}
type ListTransactionsXmlRpcResStruct struct {
	Transactions string
	Error        string
}

func (gts *GeneralTokenService) ListTransactions(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
	res := model.ListTransactionsRpcRes{}

	token := strings.ToLower(gts.Symbol)

	rpcReq := util.GenerateRpcReq(rpcConfig, strconv.Itoa(limit), token, "")
	client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcRes := ListTransactionsXmlRpcRes{}

	err := client.XmlRpcCall(gts.ParentSymbol+"Rpc.ListTransactions", &rpcReq, &rpcRes)

	if err != nil {
		return &res, err

	} else if rpcRes.Content.Error != "" {
		return &res, errors.New(rpcRes.Content.Error)

	} else if rpcRes.Content.Transactions == "" {
		return &res, errors.New("Unexpected error occured in Node.")
	}

	if err = json.Unmarshal([]byte(rpcRes.Content.Transactions), &res.Transactions); err != nil {
		return nil, err
	}

	return &res, nil
}
