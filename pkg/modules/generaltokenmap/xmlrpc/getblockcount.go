package xmlrpc

import (
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetBlockCountXmlRpcRes struct {
	Content GetBlockCountXmlRpcResStruct
}
type GetBlockCountXmlRpcResStruct struct {
	Blocks string
	Error  string
}

func (gts *GeneralTokenMapService) GetBlockCount(rpcConfig rc.RpcConfig) (res *model.GetBlockCountRpcRes, err error) {
	res = &model.GetBlockCountRpcRes{Blocks: "0"}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gts.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetBlockCountRpcRes{}, err
	}

	args, err := gts.onlyAuthArgs(rpcConfig, rpcMethod)
	if err != nil {
		return &model.GetBlockCountRpcRes{}, err
	}

	rpcReq := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gts.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.GetBlockCountRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &rpcReq, resFieldMap, res); err != nil {
		return &model.GetBlockCountRpcRes{}, err
	}

	if res.Error != "" {
		return res, errors.New(res.Error)

	} else if res.Blocks == "0" || res.Blocks == "" {
		return res, errors.New("Unexpected error occured in Node.")
	}

	return res, nil
}
