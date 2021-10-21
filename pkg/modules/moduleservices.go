package modules

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
    rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
    rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"

    generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
    generalmapxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generalmap/xmlrpc"
    generaltokenxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltoken/xmlrpc"
    generaltokenmapxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltokenmap/xmlrpc"
)

type ModuleServiceMap map[string]modules_m.ModuleService

func NewModuleServices(
    healthCheckRepo hc.HealthCheckRepository, 
    systemConfigRepo sc.SystemConfigRepository,
    rpcMethodRepo rm.Repository,
    rpcRequestRepo rrq.Repository,
    rpcResponseRepo rrs.Repository,
) *ModuleServiceMap {
    ModuleServices := make(ModuleServiceMap)
    
    // general token modules
    generalTokenModules := map[string]string{
        "THETA": "THETA",
        "TFUEL": "THETA",
        "ZIL": "ZIL",
        "TRX": "TRX",
    }
    for TOKENSYMBOL, PARENTSYMBOL := range generalTokenModules {
        ModuleServices[TOKENSYMBOL] = generaltokenxmlrpc.NewGeneralTokenService(PARENTSYMBOL, TOKENSYMBOL, healthCheckRepo, systemConfigRepo,)
    }

    // general modules
    generalModules := []string{"ALGO", "CKB", "EGLD", "FIL", "HIVE", "XTZ", "DGB", "QTUM", "HBAR"}
    for _, SYMBOL := range generalModules {
        ModuleServices[SYMBOL] = generalxmlrpc.NewGeneralService(SYMBOL, healthCheckRepo, systemConfigRepo)
    }

    // general modules using map
    generalMapModuls := []string{"BTC", "ETH"}
    for _, SYMBOL := range generalMapModuls {
        ModuleServices[SYMBOL] = generalmapxmlrpc.NewGeneralMapService(SYMBOL, healthCheckRepo, systemConfigRepo, rpcMethodRepo, rpcRequestRepo, rpcResponseRepo)
    }

    generalTokenMapModules := map[string]string{
        "ABYSS": "ETH",
    }
    for TOKENSYMBOL, PARENTSYMBOL := range generalTokenMapModules {
        ModuleServices[TOKENSYMBOL] = generaltokenmapxmlrpc.NewGeneralTokenMapService(PARENTSYMBOL, TOKENSYMBOL, healthCheckRepo, systemConfigRepo, rpcMethodRepo, rpcRequestRepo, rpcResponseRepo)
    }

    return &ModuleServices
}


