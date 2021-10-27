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
    tokenType := vars["token_type"]
    isGetAll := symbol == ""

    RES := make(GetBalanceHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - GetBalanceHandler For all symbols, Requesting ...", req) 
    } else {
        logger.InfoLog(" - GetBalanceHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req) 
    }

    gbcs.InvokeGetBalance(&RES, symbol, tokenType)

    // handle success response
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (gbcs *GetBalanceService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol, tokenType string) {
    rpcConfigCount := 0
    resChannel := make(chan GetBalanceRes)

    for _, curr := range config.CURRRPC {
        SYMBOL := strings.ToUpper(curr.Config.Symbol)
        TOKENTYPE := strings.ToUpper(curr.Config.Symbol)

        // if symbol is defined, only get for that symbol
        if symbol != "" && strings.ToUpper(symbol) != SYMBOL && strings.ToUpper(tokenType) != TOKENTYPE { continue }

        for _, rpcConfig := range curr.RpcConfigs {
            rpcConfigCount++
            _RES := GetBalanceRes{
                RpcConfig: RpcConfigResDetail{ 
                    RpcConfigId             : rpcConfig.Id,
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
                    logger.ErrorLog(" - InvokeGetBalance stas.moduleServices.GetModule err: "+err.Error())
                    _RES.Error = err.Error()
                    return
                }

                rpcRes, err := module.GetBalance(rpcConfig)
                if err != nil { 
                    logger.ErrorLog(" -- InvokeGetBalance (*gbcs.moduleServices)[SYMBOL].GetBalance(rpcConfig) Error: "+err.Error())
                    _RES.Error = rpcRes.Error

                } else {
                    logger.Log(" -- InvokeGetBalance Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance) 
                    _RES.Balance = rpcRes.Balance
                    _RES.Error   = rpcRes.Error
                }

                resChannel <- _RES

            }(curr.Config, rpcConfig)
        }
    }

    i := 0
    for res := range resChannel {
        i++
        (*RES)[res.RpcConfig.Symbol] = append((*RES)[res.RpcConfig.Symbol], res)
        if i >= rpcConfigCount { close(resChannel) }
    }
}


