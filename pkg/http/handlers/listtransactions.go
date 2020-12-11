package handlers

import (
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

type ListTransactionsHandlerResponseMap map[string][]ListTransactionsRes

type ListTransactionsService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewListTransactionsService(moduleServices *modules.ModuleServiceMap) *ListTransactionsService {
    return &ListTransactionsService{
        moduleServices,
    }
}

func (lts *ListTransactionsService) ListTransactionsHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol := vars["symbol"]
    limit  := vars["limit"]
    isGetAll := symbol == ""

    RES := make(ListTransactionsHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - ListTransactionsHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - ListTransactionsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    limitInt, _ := strconv.Atoi(limit)
    lts.InvokeListTransactions(&RES, symbol, limitInt)

    // handle success response
    logger.InfoLog(" - ListTransactionsHandler Success. Symbol: "+symbol, req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (lts *ListTransactionsService) InvokeListTransactions(RES *ListTransactionsHandlerResponseMap, symbol string, limit int) {
    rpcConfigCount := 0
    resChannel := make(chan ListTransactionsRes)

    for SYMBOL, currConfig := range config.CURR {
        SYMBOL = strings.ToUpper(SYMBOL)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            rpcConfigCount++

            _RES := ListTransactionsRes{
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
                rpcRes, err := (*lts.moduleServices)[SYMBOL].ListTransactions(rpcConfig, limit)
                if err != nil { 
                    logger.ErrorLog(" - ListTransactionsHandler (*lts.moduleServices)[SYMBOL].ListTransactions(rpcConfig, limit) Error: "+err.Error())
                    _RES.Error = rpcRes.Error

                } else {
                    logger.Log(" - InvokeListTransactions Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host) 
                    _RES.Transactions = rpcRes.Transactions
                    _RES.Error        = rpcRes.Error
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


