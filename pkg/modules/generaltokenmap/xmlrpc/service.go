package xmlrpc

import (
	"github.com/btcid/wallet-services-backend-go/cmd/config"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GeneralTokenMapService struct {
	ParentSymbol     string
	Symbol           string
	healthCheckRepo  hc.HealthCheckRepository
	systemConfigRepo sc.SystemConfigRepository
	rpcMethodRepo    rm.Repository
	rpcRequestRepo   rrq.Repository
	rpcResponseRepo  rrs.Repository
}

func (gts *GeneralTokenMapService) GetSymbol() string {
	return gts.Symbol
}

func (gts *GeneralTokenMapService) GetParentSymbol() string {
	return gts.ParentSymbol
}

func (gts *GeneralTokenMapService) GetHealthCheckRepo() hc.HealthCheckRepository {
	return gts.healthCheckRepo
}

func NewGeneralTokenMapService(
	parentSymbol string,
	symbol string,
	healthCheckRepo hc.HealthCheckRepository,
	systemConfigRepo sc.SystemConfigRepository,
	rpcMethodRepo rm.Repository,
	rpcRequestRepo rrq.Repository,
	rpcResponsRepo rrs.Repository,
) *GeneralTokenMapService {
	return &GeneralTokenMapService{
		parentSymbol,
		symbol,
		healthCheckRepo,
		systemConfigRepo,
		rpcMethodRepo,
		rpcRequestRepo,
		rpcResponsRepo,
	}
}

func (gs *GeneralTokenMapService) onlyAuthArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod) (args []string, err error) {
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
