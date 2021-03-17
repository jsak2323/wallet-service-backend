package modules

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"

    generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
    generaltokenxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltoken/xmlrpc"

    "github.com/btcid/wallet-services-backend-go/pkg/modules/btc/btcxmlrpc"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/eth/ethxmlrpc"
)

type ModuleServiceMap map[string]modules_m.ModuleService

func NewModuleServices(healthCheckRepo hc.HealthCheckRepository) *ModuleServiceMap {
    ModuleServices := make(ModuleServiceMap)

    // unique modules
    ModuleServices["BTC"] = btcxmlrpc.NewBtcService(healthCheckRepo)
    ModuleServices["ETH"] = ethxmlrpc.NewEthService(healthCheckRepo)

    // theta modules
    ModuleServices["THETA"] = generaltokenxmlrpc.NewGeneralTokenService("THETA", "THETA", healthCheckRepo)
    ModuleServices["TFUEL"] = generaltokenxmlrpc.NewGeneralTokenService("THETA", "TFUEL", healthCheckRepo)

    // tron modules
    ModuleServices["TRX"] = generaltokenxmlrpc.NewGeneralTokenService("TRX", "TRX", healthCheckRepo)

    // general modules
    generalModules := []string{"ALGO", "CKB", "EGLD", "FIL", "HIVE", "XTZ", "ZIL", "DGB", "QTUM", "HBAR"}
    for _, SYMBOL := range generalModules {
        ModuleServices[SYMBOL] = generalxmlrpc.NewGeneralService(SYMBOL, healthCheckRepo)
    }

    return &ModuleServices
}


