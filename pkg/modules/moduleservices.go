package modules

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"

    generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
    generaltokenxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltoken/xmlrpc"

    "github.com/btcid/wallet-services-backend-go/pkg/modules/btc/btcxmlrpc"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/eth/ethxmlrpc"
)

type ModuleServiceMap map[string]modules_m.ModuleService

func NewModuleServices(
    healthCheckRepo hc.HealthCheckRepository, 
    systemConfigRepo sc.SystemConfigRepository,
) *ModuleServiceMap {
    ModuleServices := make(ModuleServiceMap)

    // unique modules
    ModuleServices["BTC"] = btcxmlrpc.NewBtcService(healthCheckRepo, systemConfigRepo)
    ModuleServices["ETH"] = ethxmlrpc.NewEthService(healthCheckRepo, systemConfigRepo)

    // theta modules
    ModuleServices["THETA"] = generaltokenxmlrpc.NewGeneralTokenService("THETA", "THETA", healthCheckRepo, systemConfigRepo)
    ModuleServices["TFUEL"] = generaltokenxmlrpc.NewGeneralTokenService("THETA", "TFUEL", healthCheckRepo, systemConfigRepo)

    // tron modules
    ModuleServices["TRX"] = generaltokenxmlrpc.NewGeneralTokenService("TRX", "TRX", healthCheckRepo, systemConfigRepo)

    // general modules
    generalModules := []string{"ALGO", "CKB", "EGLD", "FIL", "HIVE", "XTZ", "ZIL", "DGB", "QTUM", "HBAR"}
    for _, SYMBOL := range generalModules {
        ModuleServices[SYMBOL] = generalxmlrpc.NewGeneralService(SYMBOL, healthCheckRepo, systemConfigRepo)
    }

    return &ModuleServices
}


