package handlers

import (
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

type GetBlockCountHandlerResponseMap map[string]map[string][]GetBlockCountRes

type GetBlockCountService struct {
    moduleServices   *modules.ModuleServiceMap
    systemConfigRepo sc.SystemConfigRepository
}

func NewGetBlockCountService(
    moduleServices *modules.ModuleServiceMap, 
    systemConfigRepo sc.SystemConfigRepository,
) *GetBlockCountService {
    return &GetBlockCountService{
        moduleServices,
        systemConfigRepo,
    }
}

func (gbcs *GetBlockCountService) GetBlockCountHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    tokenType := vars["token_type"]
    isGetAll := symbol == ""

    RES := make(GetBlockCountHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBlockCountHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBlockCount(&RES, symbol, tokenType)
    
    // handle success response
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBlockCountService) InvokeGetBlockCount(RES *GetBlockCountHandlerResponseMap, symbol, tokenType string) {
    rpcConfigCount := 0
    resChannel := make(chan GetBlockCountRes)

    maintenanceList, err := GetMaintenanceList(gbcs.systemConfigRepo)
    if err != nil { logger.ErrorLog(" - InvokeGetBlockCount h.GetMaintenanceList err: "+err.Error()) }

    for _, currRpc := range config.CURRRPC {
        SYMBOL := strings.ToUpper(currRpc.Config.Symbol)
        TOKENTYPE := strings.ToUpper(currRpc.Config.TokenType)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL && strings.ToUpper(tokenType) != TOKENTYPE { continue }

        // if not parent coin, skip
        if currRpc.Config.ParentSymbol != cc.MainTokenType { continue }

        // if maintenance, skip
        if maintenanceList[SYMBOL] { continue }

        for _, rpcConfig := range currRpc.RpcConfigs {
            rpcConfigCount++

            _RES := GetBlockCountRes{
                RpcConfig: RpcConfigResDetail{
                    RpcConfigId             : rpcConfig.Id,
                    Symbol                  : SYMBOL,
                    Name                    : rpcConfig.Name,
                    Host                    : rpcConfig.Host,
                    Type                    : rpcConfig.Type,
                    NodeVersion             : rpcConfig.NodeVersion,
                    NodeLastUpdated         : rpcConfig.NodeLastUpdated,
                    IsHealthCheckEnabled    : rpcConfig.IsHealthCheckEnabled,
                },
            }

            // execute concurrent rpc calls
            go func(currencyConfig cc.CurrencyConfig, rpcConfig rc.RpcConfig) {
                module, err := gbcs.moduleServices.GetModule(currencyConfig.Id)
                if err != nil {
                    logger.ErrorLog(" - InvokeGetBlockCount stas.moduleServices.GetModule err: "+err.Error())
                    _RES.Error = err.Error()
                    return
                }
                
                rpcRes, err := module.GetBlockCount(rpcConfig)
                if err != nil { 
                    logger.Log(" - InvokeGetBlockCount (*gbcs.moduleServices)["+SYMBOL+"].GetBlockCount(rpcConfig) Error: "+err.Error())
                    _RES.Error = rpcRes.Error

                } else {
                    logger.Log(" - InvokeGetBlockCount Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Blocks: "+rpcRes.Blocks) 
                    _RES.Blocks = rpcRes.Blocks
                    _RES.Error  = rpcRes.Error
                }

                resChannel <- _RES
                
            }(currRpc.Config, rpcConfig)
        }
    }

    i := 0
    for res := range resChannel {
        i++
        _, ok := (*RES)[res.RpcConfig.Symbol]
        if !ok { (*RES)[res.RpcConfig.Symbol] = make(map[string][]GetBlockCountRes) }

        (*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}


