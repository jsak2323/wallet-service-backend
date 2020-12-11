package modules

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/btc/btcxmlrpc"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/eth/ethxmlrpc"
)

type ModuleServiceMap map[string]modules_m.ModuleService

func NewModuleServices(healthCheckRepo hc.HealthCheckRepository) *ModuleServiceMap {

    generalModules := []string{"ALGO", "CKB", "EGLD", "FIL", "HIVE", "XTZ", "ZIL", "DGB"}

    ModuleServices := make(ModuleServiceMap)

    // unique modules
    ModuleServices["BTC"] = btcxmlrpc.NewBtcService(healthCheckRepo)
    ModuleServices["ETH"] = ethxmlrpc.NewEthService(healthCheckRepo)

    // general modules
    for _, SYMBOL := range generalModules {
        ModuleServices[SYMBOL] = generalxmlrpc.NewGeneralService(SYMBOL, healthCheckRepo)
    }

    return &ModuleServices
}


