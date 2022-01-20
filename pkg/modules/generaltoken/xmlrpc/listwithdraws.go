package xmlrpc

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type ListWithdrawsXmlRpcRes struct {
	Content ListWithdrawsXmlRpcResStruct
}
type ListWithdrawsXmlRpcResStruct struct {
	Withdraws string
	Error     string
}

func (gts *GeneralTokenService) ListWithdraws(rpcConfig rc.RpcConfig, limit int) (*model.ListWithdrawsRpcRes, error) {
	res := model.ListWithdrawsRpcRes{}

	token := strings.ToLower(gts.Symbol)

	rpcReq := util.GenerateRpcReq(rpcConfig, strconv.Itoa(limit), token, "")
	client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcRes := ListWithdrawsXmlRpcRes{}

	err := client.XmlRpcCall(gts.ParentSymbol+"Rpc.ListWithdraws", &rpcReq, &rpcRes)

	if err != nil {
		return &res, err

	} else if rpcRes.Content.Error != "" {
		return &res, errors.New(rpcRes.Content.Error)

	} else if rpcRes.Content.Withdraws == "" {
		return &res, errors.New("Unexpected error occured in Node.")
	}

	if err = json.Unmarshal([]byte(rpcRes.Content.Withdraws), &res.Withdraws); err != nil {
		return nil, err
	}

	return &res, nil
}
