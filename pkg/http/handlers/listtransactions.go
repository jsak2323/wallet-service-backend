package handlers

import (
    "sync"
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

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - ListTransactionsHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (lts *ListTransactionsService) InvokeListTransactions(RES *ListTransactionsHandlerResponseMap, symbol string, limit int) {
    var wg sync.WaitGroup
    rpcConfigCount := 0
    resChannel := make(chan ListTransactionsRes)

    for confKey, currConfig := range config.CURR {
        confKey = strings.ToUpper(confKey)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != confKey { continue }

        for _, rpcConfig := range currConfig.RpcConfigs {
            wg.Add(1)
            rpcConfigCount++
            wg.Done()

            go func(confKey string, rpcConfig rc.RpcConfig) {
                rpcRes, err := (*lts.moduleServices)[confKey].ListTransactions(rpcConfig, limit)
                if err != nil { 
                    logger.ErrorLog(" - ListTransactionsHandler (*lts.moduleServices)[confKey].ListTransactions(rpcConfig, limit) Error: "+err.Error())
                }

                logger.Log(" - InvokeListTransactions Symbol: "+confKey+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance) 
                resChannel <- ListTransactionsRes{
                    RpcConfig: RpcConfigResDetail{
                        RpcConfigId         : rpcConfig.Id,
                        Symbol              : confKey,
                        Name                : rpcConfig.Name,
                        Host                : rpcConfig.Host,
                        Type                : rpcConfig.Type,
                        NodeVersion         : rpcConfig.NodeVersion,
                        NodeLastUpdated     : rpcConfig.NodeLastUpdated,
                    },
                    Transactions: rpcRes.Transactions,
                }
            }(confKey, rpcConfig)
        }
    }

    wg.Wait()
    i := 0
    for res := range resChannel {
        i++
        (*RES)[res.Symbol] = append((*RES)[res.Symbol], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}
