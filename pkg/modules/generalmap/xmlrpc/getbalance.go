package xmlrpc

import (
	"errors"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
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

func (gms *GeneralMapService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
	res := &model.GetBalanceRpcRes{Balance: "0"}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gms.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	runtimeParams := map[string]string{
		"token": strings.ToLower(gms.Symbol),
	}

	args, err := util.GetRpcRequestArgs(rpcConfig, rpcMethod, rpcRequests, runtimeParams)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gms.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &req, resFieldMap, res); err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	if res.Error != "" {
		return &model.GetBalanceRpcRes{}, errors.New(res.Error)
	} else if res.Balance == "0" || res.Balance == "" {
		return &model.GetBalanceRpcRes{}, errors.New("Unexpected error occured in Node.")
	}

	return res, nil
}
