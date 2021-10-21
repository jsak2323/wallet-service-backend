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

type GetBalanceXmlRpcRes struct {
	Content GetBalanceXmlRpcResStruct
}
type GetBalanceXmlRpcResStruct struct {
	Balance string
	Error   string
}

func (gts *GeneralTokenMapService) GetBalance(rpcConfig rc.RpcConfig) (res *model.GetBalanceRpcRes, err error) {
	res = &model.GetBalanceRpcRes{Balance: "0"}

	client := util.NewXmlRpcMapClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	rpcMethod, err := config.GetRpcMethod(gts.rpcMethodRepo, rpcConfig.Id, rm.TypeGetBalance)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	token := strings.ToLower(gts.Symbol)

	args, err := gts.getBalanceArgs(rpcConfig, rpcMethod, token)
	if err != nil {
		return &model.GetBalanceRpcRes{}, err
	}

	req := util.GenerateRpcMapRequest(args)

	resFieldMap, err := config.GetRpcResponseMap(gts.rpcResponseRepo, rpcMethod.Id)
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

func (gs *GeneralTokenMapService) getBalanceArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, token string) (args []string, err error) {
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
