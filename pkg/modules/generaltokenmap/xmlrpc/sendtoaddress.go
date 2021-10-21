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

func (gts *GeneralTokenMapService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (res *model.SendToAddressRpcRes, err error) {
	token := strings.ToLower(gts.Symbol)

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gts.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	args, err := gts.sendToAddressArgs(rpcConfig, rpcMethod, token, amountInDecimal, address, memo)
	if err != nil {
		return &model.SendToAddressRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gts.rpcResponseRepo, rpcMethod.Id)
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

func (gts *GeneralTokenMapService) sendToAddressArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, token, amountInDecimal, address, memo string) (args []string, err error) {
	args = make([]string, rpcMethod.NumOfArgs)

	hashkey, nonce := util.GenerateHashkey(rpcConfig.Password, rpcConfig.Hashkey)

	rpcRequests, err := config.GetRpcRequestMap(gts.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return []string{}, err
	}

	for _, rpcRequest := range rpcRequests {
		if rpcRequest.Source == rrq.SourceRuntime {
			switch rpcRequest.ArgName {
			case rrq.ArgRpcUser:
				args[rpcRequest.ArgOrder] = rpcConfig.User
			case rrq.ArgHashkey:
				args[rpcRequest.ArgOrder] = hashkey
			case rrq.ArgNonce:
				args[rpcRequest.ArgOrder] = nonce
			case rrq.ArgToken:
				args[rpcRequest.ArgOrder] = token
			case rrq.ArgAddress:
				args[rpcRequest.ArgOrder] = address
			case rrq.ArgAmount:
				args[rpcRequest.ArgOrder] = amountInDecimal
			case rrq.ArgMemo:
				args[rpcRequest.ArgOrder] = memo
			default:
				return []string{}, model.InvalidRpcRequestConfig(rpcRequest.ArgName, rpcMethod.Name)
			}
		}

		if rpcRequest.Source == rrq.SourceConfig {
			args = append(args, rpcRequest.Value)
		}
	}

	return args, nil
}
