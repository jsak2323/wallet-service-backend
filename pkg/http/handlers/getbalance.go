package handlers

import (
    "fmt"
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules"
    "github.com/btcid/wallet-services-backend/cmd/config"
)

type GetBalanceHandlerResponseMap map[string][]GetBalanceRes

type GetBalanceService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetBalanceService(moduleServices *modules.ModuleServiceMap) *GetBalanceService {
    return &GetBalanceService{
        moduleServices,
    }
}

func (gbcs *GetBalanceService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    isGetAll := symbol == ""

    RES := make(GetBalanceHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBalanceHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBalanceHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBalance(&RES, symbol)

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBalanceService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol string) {
    rpcConfigCount := 0
    resChannel := make(chan GetBalanceRes)

    for SYMBOL, currConfig := range config.CURR {
        SYMBOL = strings.ToUpper(SYMBOL)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            rpcConfigCount++

            _RES := GetBalanceRes{
                RpcConfig: RpcConfigResDetail{ 
                    RpcConfigId         : rpcConfig.Id,
                    Symbol              : SYMBOL,
                    Name                : rpcConfig.Name,
                    Host                : rpcConfig.Host,
                    Type                : rpcConfig.Type,
                    NodeVersion         : rpcConfig.NodeVersion,
                    NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                },
            }

            go func(SYMBOL string, rpcConfig rc.RpcConfig) {
                fmt.Println(" - SYMBOL: ", SYMBOL)
                rpcRes, err := (*gbcs.moduleServices)[SYMBOL].GetBalance(rpcConfig)
                if err != nil { 
                    logger.ErrorLog(" -- InvokeGetBalance (*gbcs.moduleServices)[SYMBOL].GetBalance(rpcConfig) Error: "+err.Error())
                    _RES.Error = rpcRes.Error

                } else {
                    logger.Log(" -- InvokeGetBalance Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance) 
                    _RES.Balance = rpcRes.Balance
                    _RES.Error   = rpcRes.Error
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


