package modules

import (
    "errors"
    "strconv"
    
	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"

	generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
	generalmapxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generalmap/xmlrpc"
	generaltokenxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltoken/xmlrpc"
)

type ModuleServiceMap map[int]modules_m.ModuleService

func NewModuleServices(
    healthCheckRepo hc.Repository, 
    systemConfigRepo sc.Repository,
    rpcMethodRepo rm.Repository,
    rpcRequestRepo rrq.Repository,
    rpcResponseRepo rrs.Repository,
) *ModuleServiceMap {
    ModuleServices := make(ModuleServiceMap)
    
    for currencyConfigId, curr := range config.CURRRPC {
        switch curr.Config.ModuleType {
        case cc.ModuleTypeGeneral:
            ModuleServices[currencyConfigId] = generalxmlrpc.NewGeneralService(curr.Config, healthCheckRepo, systemConfigRepo)
        case cc.ModuleTypeGeneralToken:
            ModuleServices[currencyConfigId] = generaltokenxmlrpc.NewGeneralTokenService(curr.Config, healthCheckRepo, systemConfigRepo)
        case cc.ModuleTypeGeneralMap:
            ModuleServices[currencyConfigId] = generalmapxmlrpc.NewGeneralMapService(curr.Config, healthCheckRepo, systemConfigRepo, rpcMethodRepo, rpcRequestRepo, rpcResponseRepo)
        }
    }

    return &ModuleServices
}

func(msm *ModuleServiceMap) GetModule(currencyConfigId int) (modules_m.ModuleService, error) {
    module, ok := (*msm)[currencyConfigId]
    if !ok {
        return nil, errors.New("module not implemented for currencyConfigId: "+strconv.Itoa(currencyConfigId))
    }

    return module, nil
}


