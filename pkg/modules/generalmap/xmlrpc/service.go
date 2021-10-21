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
)

type GeneralMapService struct {
	Symbol           string
	healthCheckRepo  hc.HealthCheckRepository
	systemConfigRepo sc.SystemConfigRepository
	rpcMethodRepo    rm.Repository
	rpcRequestRepo   rrq.Repository
	rpcResponseRepo  rrs.Repository
}

func (gs *GeneralMapService) GetSymbol() string {
	return gs.Symbol
}

func (gs *GeneralMapService) GetHealthCheckRepo() hc.HealthCheckRepository {
	return gs.healthCheckRepo
}

func NewGeneralMapService(
	symbol string,
	healthCheckRepo hc.HealthCheckRepository,
	systemConfigRepo sc.SystemConfigRepository,
	rpcMethodRepo rm.Repository,
	rpcRequestRepo rrq.Repository,
	rpcResponsRepo rrs.Repository,
) *GeneralMapService {
	return &GeneralMapService{
		symbol,
		healthCheckRepo,
		systemConfigRepo,
		rpcMethodRepo,
		rpcRequestRepo,
		rpcResponsRepo,
	}
}

func (gs *GeneralMapService) onlyAuthArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod) (args []string, err error) {
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
				// log
			}
		}

		if rpcRequest.Source == rrq.SourceConfig {
			args = append(args, rpcRequest.Value)
		}
	}

	return args, nil
}
