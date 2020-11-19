package handlers

import (
    // "sync"
    // "strconv"
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    // rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
    "github.com/btcid/wallet-services-backend/pkg/modules"
    // "github.com/btcid/wallet-services-backend/cmd/config"
)

type SendToAddressService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewSendToAddressService(moduleServices *modules.ModuleServiceMap) *SendToAddressService {
    return &SendToAddressService{
        moduleServices,
    }
}

func (stas *SendToAddressService) SendToAddressHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol          := vars["symbol"]
    address         := vars["address"]
    amountInDecimal := vars["amount"]

    SYMBOL := strings.ToUpper(symbol)

    logger.InfoLog(" - SendToAddressHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "sender")
    if err != nil {
        logger.ErrorLog(" - SendToAddressHandler util.GetRpcConfigByType(SYMBOL, \"sender\") err: "+err.Error())
        return
    }

    rpcRes, err := (*stas.moduleServices)[SYMBOL].SendToAddress(rpcConfig, address, amountInDecimal)
    if err != nil { 
        logger.ErrorLog(" - SendToAddressHandler (*stas.moduleServices)[strings.ToUpper(symbol)].SendToAddress(rpcConfig, address, amountInDecimal) address:"+address+", amount: "+amountInDecimal+", Error: "+err.Error())
        return
    }

    RES := SendToAddressRes{
        RpcConfig: RpcConfigResDetail{ 
            RpcConfigId         : rpcConfig.Id,
            Symbol              : SYMBOL,
            Name                : rpcConfig.Name,
            Host                : rpcConfig.Host,
            Type                : rpcConfig.Type,
            NodeVersion         : rpcConfig.NodeVersion,
            NodeLastUpdated     : rpcConfig.NodeLastUpdated,
        },
        TxHash  : rpcRes.TxHash,
        Error   : rpcRes.Error,
    }

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - SendToAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}


