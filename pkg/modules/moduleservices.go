package modules

import(
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
    rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
    rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    modules_m "github.com/btcid/wallet-services-backend-go/pkg/modules/model"

    generalxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/general/xmlrpc"
    generaltokenxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generaltoken/xmlrpc"
    generalmapxmlrpc "github.com/btcid/wallet-services-backend-go/pkg/modules/generalmap/xmlrpc"
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

    generalTokenMapModules := map[string]string{
        "ABYSS": "ETH",
        "BTC": "BTC",
        "ETH": "ETH",
    }
    for SYMBOL, PARENTSYMBOL := range generalTokenMapModules {
        ModuleServices[SYMBOL] = generalmapxmlrpc.NewGeneralMapService(PARENTSYMBOL, SYMBOL, healthCheckRepo, systemConfigRepo, rpcMethodRepo, rpcRequestRepo, rpcResponseRepo)
    }

    return &ModuleServices
}


