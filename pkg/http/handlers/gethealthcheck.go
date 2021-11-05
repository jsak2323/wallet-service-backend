package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

type GetHealthCheckHandlerResponseMap map[string]map[string][]GetHealthCheckRes // map by symbol, token_type

type GetHealthCheckService struct {
    moduleServices   *modules.ModuleServiceMap
    healthCheckRepo  hc.HealthCheckRepository
    systemConfigRepo sc.SystemConfigRepository
}

func NewGetHealthCheckService(
    moduleServices *modules.ModuleServiceMap, 
    healthCheckRepo hc.HealthCheckRepository,
    systemConfigRepo sc.SystemConfigRepository,
) *GetHealthCheckService {
    return &GetHealthCheckService{
        moduleServices,
        healthCheckRepo,
        systemConfigRepo,
    }
}

func (ghcs *GetHealthCheckService) GetHealthCheckHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    tokenType := vars["token_type"]
    isGetAll := symbol == ""

    RES := make(GetHealthCheckHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetHealthCheckHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetHealthCheckHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    ghcs.InvokeGetHealthCheck(&RES, symbol, tokenType)
    
    // handle success response
    logger.InfoLog(" - GetHealthCheckHandler Success.", req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (ghcs *GetHealthCheckService) InvokeGetHealthCheck(RES *GetHealthCheckHandlerResponseMap, symbol, tokenType string) {
    // get maintenance list
    maintenanceList, err := GetMaintenanceList(ghcs.systemConfigRepo)
    if err != nil { logger.ErrorLog(" - InvokeGetHealthCheck GetMaintenanceList err: "+err.Error()) }

    // fetch healthcheck data from db
    if symbol != "" { // get by rpc config id
        SYMBOL := strings.ToUpper(symbol)
        TOKENTYPE := strings.ToUpper(tokenType)

        currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, TOKENTYPE)
        if err != nil {
            logger.ErrorLog(" - InvokeGetHealthCheck "+SYMBOL+" config.GetCurrencyBySymbol err: "+err.Error())
            return
        }

        for _, rpcConfig := range config.CURRRPC[currencyConfig.Id].RpcConfigs {
            _RES := GetHealthCheckRes{}

            if rpcConfig.Id == 0 { // currency not found
                _RES.Error = "Invalid currency."
                return 
            } 
            _RES.RpcConfig = RpcConfigResDetail{
                RpcConfigId          : rpcConfig.Id,
                Symbol               : SYMBOL,
                TokenType            : TOKENTYPE,
                Name                 : rpcConfig.Name,
                Host                 : rpcConfig.Host,
                Type                 : rpcConfig.Type,
                NodeVersion          : rpcConfig.NodeVersion,
                NodeLastUpdated      : rpcConfig.NodeLastUpdated,
                IsHealthCheckEnabled : rpcConfig.IsHealthCheckEnabled,
            }

            healthCheck, err := ghcs.healthCheckRepo.GetByRpcConfigId(rpcConfig.Id)
            if err != nil {
                logger.ErrorLog(" - InvokeGetHealthCheck "+SYMBOL+" GetByRpcConfigId(rpcConfig.Id) err: "+err.Error())
                _RES.Error = err.Error()
                return
            }

            _RES.HealthCheck   = healthCheck
            _RES.IsMaintenance = maintenanceList[SYMBOL]
            
            _, ok := (*RES)[SYMBOL]
            if !ok { (*RES)[SYMBOL] = make(map[string][]GetHealthCheckRes) }
            (*RES)[SYMBOL][TOKENTYPE] = append((*RES)[SYMBOL][TOKENTYPE], _RES)
        }

    } else { // get all

        for _, curr := range config.CURRRPC {
            for _, rpcConfig := range curr.RpcConfigs {

                healthCheck, err := ghcs.healthCheckRepo.GetByRpcConfigId(rpcConfig.Id)
                if err != nil {
                    logger.ErrorLog(" - InvokeGetHealthCheck GetAllWithRpcConfig err: "+err.Error())
                    return
                }

                _RES := GetHealthCheckRes{
                    RpcConfig: RpcConfigResDetail{
                        RpcConfigId          : rpcConfig.Id,
                        Symbol               : curr.Config.Symbol,
                        TokenType            : curr.Config.TokenType,
                        Name                 : rpcConfig.Name,
                        Host                 : rpcConfig.Host,
                        Type                 : rpcConfig.Type,
                        NodeVersion          : rpcConfig.NodeVersion,
                        NodeLastUpdated      : rpcConfig.NodeLastUpdated,
                        IsHealthCheckEnabled : rpcConfig.IsHealthCheckEnabled,
                    },
                    HealthCheck: hc.HealthCheck{
                        Id           : healthCheck.Id,
                        RpcConfigId  : rpcConfig.Id,
                        BlockCount   : healthCheck.BlockCount,
                        BlockDiff    : healthCheck.BlockDiff,
                        IsHealthy    : healthCheck.IsHealthy,
                        LastUpdated  : healthCheck.LastUpdated,
                    },
                    IsMaintenance: maintenanceList[curr.Config.Symbol],
                }

                _, ok := (*RES)[curr.Config.Symbol]
                if !ok { (*RES)[curr.Config.Symbol] = make(map[string][]GetHealthCheckRes) }
                (*RES)[curr.Config.Symbol][curr.Config.TokenType] = append((*RES)[curr.Config.Symbol][curr.Config.TokenType], _RES)
            }
        }
    }
}


