package handlers

import (
    "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
)

type ListTransactionsHandlerResponseMap map[string]map[string][]ListTransactionsRes // map by symbol, token_type

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
    tokenType := vars["token_type"]
    limit  := vars["limit"]
    isGetAll := symbol == ""

    RES := make(ListTransactionsHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - ListTransactionsHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - ListTransactionsHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    limitInt, _ := strconv.Atoi(limit)
    lts.InvokeListTransactions(&RES, symbol, tokenType, limitInt)

    // handle success response
    logger.InfoLog(" - ListTransactionsHandler Success. Symbol: "+symbol, req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (lts *ListTransactionsService) InvokeListTransactions(RES *ListTransactionsHandlerResponseMap, symbol, tokenType string, limit int) {
    rpcConfigCount := 0
    resChannel := make(chan ListTransactionsRes)

    for _, curr := range config.CURRRPC {
        SYMBOL := strings.ToUpper(curr.Config.Symbol)
        TOKENTYPE := strings.ToUpper(curr.Config.Symbol)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL {continue}
        if tokenType != "" && strings.ToUpper(tokenType) != TOKENTYPE { continue }

        for _, rpcConfig := range curr.RpcConfigs {
            rpcConfigCount++

            _RES := ListTransactionsRes{
                RpcConfig: RpcConfigResDetail{
                    RpcConfigId             : rpcConfig.Id,
                    Symbol                  : SYMBOL,
                    TokenType               : TOKENTYPE,
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
                module, err := lts.moduleServices.GetModule(currencyConfig.Id)
                if err != nil {
                    logger.ErrorLog(" - ListTransactionsHandler lts.moduleServices.GetModule err: "+err.Error())
                    _RES.Error = err.Error()
                    return
                }
                
                rpcRes, err := module.ListTransactions(rpcConfig, limit)
                if err != nil { 
                    logger.ErrorLog(" - ListTransactionsHandler (*lts.moduleServices)[SYMBOL].ListTransactions(rpcConfig, limit) Error: "+err.Error())
                    _RES.Error = rpcRes.Error

                } else {
                    logger.Log(" - InvokeListTransactions Symbol: "+currencyConfig.Symbol+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host) 
                    _RES.Transactions = rpcRes.Transactions
                    _RES.Error        = rpcRes.Error
                }

                resChannel <- _RES

            }(curr.Config, rpcConfig)
        }
    }

    i := 0
    for res := range resChannel {
        i++
        _, ok := (*RES)[res.RpcConfig.Symbol]
        if !ok {
            (*RES)[res.RpcConfig.Symbol] = make(map[string][]ListTransactionsRes)
        }

        (*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType] = append((*RES)[res.RpcConfig.Symbol][res.RpcConfig.TokenType], res)

        if i >= rpcConfigCount { close(resChannel) }
    }
}


