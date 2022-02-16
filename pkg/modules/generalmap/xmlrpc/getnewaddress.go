package xmlrpc

import (
	"context"
	"errors"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (gms *GeneralMapService) GetNewAddress(ctx context.Context, rpcConfig rc.RpcConfig, addressType string) (res *model.GetNewAddressRpcRes, err error) {
	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(ctx, gms.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	runtimeParams := map[string]string{
		"token": strings.ToLower(gms.Symbol),
	}

	args, err := util.GetRpcRequestArgs(rpcConfig, rpcMethod, rpcRequests, runtimeParams)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	rpcReq := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(ctx, gms.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &rpcReq, resFieldMap, res); err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	if res.Error != "" {
		return &model.GetNewAddressRpcRes{}, errors.New(res.Error)
	} else if res.Address == "" {
		return &model.GetNewAddressRpcRes{}, errors.New("Unexpected error occured in Node.")
	}

	return res, nil
}
