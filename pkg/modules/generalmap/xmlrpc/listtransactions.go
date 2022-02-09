package xmlrpc

import (
	"context"
	"strconv"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (gms *GeneralMapService) ListTransactions(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (res *model.ListTransactionsRpcRes, err error) {
	res = &model.ListTransactionsRpcRes{}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(ctx, gms.rpcMethodRepo, rpcConfig.Id, rm.TypeListTransactions)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	runtimeParams := map[string]string{
		"token": strings.ToLower(gms.Symbol),
		"limit": strconv.Itoa(limit),
	}

	args, err := util.GetRpcRequestArgs(rpcConfig, rpcMethod, rpcRequests, runtimeParams)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gms.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &req, resFieldMap, res); err != nil {
		return &model.ListTransactionsRpcRes{}, err
	}

	return res, nil
}

func (gms *GeneralMapService) listTransactionsArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, limit, token string) (args []string, err error) {
	args = make([]string, rpcMethod.NumOfArgs)

	hashkey, nonce := util.GenerateHashkey(rpcConfig.Password, rpcConfig.Hashkey)

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
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
			case rrq.ArgToken:
				args[rpcRequest.ArgOrder] = token
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
