package xmlrpc

import (
	"errors"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
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

func (gs *GeneralMapService) ListTransactions(rpcConfig rc.RpcConfig, limit int) (res *model.ListTransactionsRpcRes, err error) {
	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gs.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	args, err := gs.listTransactionsArgs(rpcConfig, rpcMethod, strconv.Itoa(limit))
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	rpcReq := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gs.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &rpcReq, resFieldMap, res); err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	if res.Error != "" {
		return &model.ListTransactionsRpcRes{}, errors.New(res.Error)
	} else if res.Transactions == "" {
		return &model.ListTransactionsRpcRes{}, errors.New("Unexpected error occured in Node.")
	}

	return res, nil
}

func (gs *GeneralMapService) listTransactionsArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, limit string) (args []string, err error) {
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
				args[rpcRequest.ArgOrder] = limit
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
