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

func (gms *GeneralMapService) ListWithdraws(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (res *model.ListWithdrawsRpcRes, err error) {
	res = &model.ListWithdrawsRpcRes{}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(ctx, gms.rpcMethodRepo, rpcConfig.Id, rm.TypeListWithdraws)
	if err != nil {
		return &model.ListWithdrawsRpcRes{}, err
	}

	rpcRequests, err := config.GetRpcRequestMap(gms.rpcRequestRepo, rpcMethod.Id)
	if err != nil {
		return &model.ListWithdrawsRpcRes{}, err
	}

	runtimeParams := map[string]string{
		"token": strings.ToLower(gms.Symbol),
		"limit": strconv.Itoa(limit),
	}

	args, err := util.GetRpcRequestArgs(rpcConfig, rpcMethod, rpcRequests, runtimeParams)
	if err != nil {
		return &model.ListWithdrawsRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gms.rpcResponseRepo, rpcMethod.Id)
	if err != nil {
		return &model.ListWithdrawsRpcRes{}, err
	}

	if err = client.XmlRpcMapCall(rpcMethod.Name, &req, resFieldMap, res); err != nil {
		return &model.ListWithdrawsRpcRes{}, err
	}

	return res, nil
}

func (gms *GeneralMapService) listWithdrawsArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, limit, token string) (args []string, err error) {
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
