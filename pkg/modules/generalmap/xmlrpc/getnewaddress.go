package xmlrpc

import (
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetNewAddressXmlRpcRes struct {
	Content GetNewAddressXmlRpcResStruct
}
type GetNewAddressXmlRpcResStruct struct {
	Address string
	Error   string
}

func (gs *GeneralMapService) GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (res *model.GetNewAddressRpcRes, err error) {
	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gs.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	args, err := gs.getNewAddressArgs(rpcConfig, rpcMethod, addressType)
	if err != nil {
		return &model.GetNewAddressRpcRes{}, err
	}

	rpcReq := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gs.rpcResponseRepo, rpcMethod.Id)
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

func (gs *GeneralMapService) getNewAddressArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, addresstype string) (args []string, err error) {
	args = make([]string, rpcMethod.NumOfArgs)

	hashkey, nonce := util.GenerateHashkey(rpcConfig.Password, rpcConfig.Hashkey)

	rpcRequests, err := config.GetRpcRequestMap(gs.rpcRequestRepo, rpcMethod.Id)
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
			case rrq.ArgAddressType:
				args[rpcRequest.ArgOrder] = addresstype
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
