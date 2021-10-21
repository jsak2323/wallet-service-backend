package xmlrpc

import (
	"errors"

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

func (gs *GeneralMapService) GetBalance(rpcConfig rc.RpcConfig) (res *model.GetBalanceRpcRes, err error) {
	res = &model.GetBalanceRpcRes{Balance: "0"}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gs.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	args, err := gs.onlyAuthArgs(rpcConfig, rpcMethod)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gs.rpcResponseRepo, rpcMethod.Id)
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
