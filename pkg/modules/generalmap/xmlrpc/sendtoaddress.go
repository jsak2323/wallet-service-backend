package xmlrpc

import (
	"errors"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (gms *GeneralMapService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
	res := &model.SendToAddressRpcRes{}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gms.rpcMethodRepo, rpcConfig.Id, rm.TypeSendToAddress)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	runtimeParams := map[string]string{
		rrq.ArgToken:   strings.ToLower(gms.Symbol),
		rrq.ArgAmount:  amountInDecimal,
		rrq.ArgAddress: address,
		rrq.ArgMemo:    memo,
	}

	args, err := util.GetRpcRequestArgs(rpcConfig, rpcMethod, rpcRequests, runtimeParams)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gms.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &req, resFieldMap, res); err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	if res.Error != "" {
		return res, errors.New(res.Error)

	} else if res.TxHash == "" {
		return res, errors.New("Unexpected error occured in Node.")
	}

	return res, nil
}
