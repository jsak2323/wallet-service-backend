package handlers

import (
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

type GetBlockCountHandlerResponseMap map[string][]GetBlockCountRes

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
    isGetAll := symbol == ""

    RES := make(GetBlockCountHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBlockCountHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBlockCountHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBlockCount(&RES, symbol)
    
    // handle success response
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBlockCountHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBlockCountService) InvokeGetBlockCount(RES *GetBlockCountHandlerResponseMap, symbol string) {
    rpcConfigCount := 0
    resChannel := make(chan GetBlockCountRes)

    maintenanceList, err := GetMaintenanceList(gbcs.systemConfigRepo)
    if err != nil { logger.ErrorLog(" - InvokeGetBlockCount h.GetMaintenanceList err: "+err.Error()) }

    for SYMBOL, currConfig := range config.CURR {
        SYMBOL = strings.ToUpper(SYMBOL)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL { continue }

        // if not parent coin, skip
        if currConfig.Config.ParentSymbol != "" { continue }

        // if maintenance, skip
        if maintenanceList[SYMBOL] { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
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
            go func(SYMBOL string, rpcConfig rc.RpcConfig) {
                module, ok := (*gbcs.moduleServices)[SYMBOL]
                if !ok {
                    logger.ErrorLog(" - InvokeGetBlockCount module not implemented symbol: "+SYMBOL)
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
                
            }(SYMBOL, rpcConfig)
        }
    }

    i := 0
    for res := range resChannel {
        i++
        (*RES)[res.RpcConfig.Symbol] = append((*RES)[res.RpcConfig.Symbol], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}


