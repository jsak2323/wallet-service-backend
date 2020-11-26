package modules

import(
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
    modules_m "github.com/btcid/wallet-services-backend/pkg/modules/model"
    generalxmlrpc "github.com/btcid/wallet-services-backend/pkg/modules/general/xmlrpc"
    "github.com/btcid/wallet-services-backend/pkg/modules/btc/btcxmlrpc"

    // "github.com/btcid/wallet-services-backend/pkg/modules/eth/ethxmlrpc"
)

type ModuleServiceMap map[string]modules_m.ModuleService

func NewModuleServices(healthCheckRepo hc.HealthCheckRepository) *ModuleServiceMap {

    ModuleServices := make(ModuleServiceMap)

    ModuleServices["BTC"] = btcxmlrpc.NewBtcService(healthCheckRepo)
    // ModuleServices["ETH"] = ethxmlrpc.NewEthService(healthCheckRepo)
    ModuleServices["ALGO"] = generalxmlrpc.NewGeneralService("ALGO", healthCheckRepo)


    return &ModuleServices
}


