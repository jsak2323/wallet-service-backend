package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
    "github.com/btcid/wallet-services-backend/pkg/modules"
)

type GetNewAddressService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewGetNewAddressService(moduleServices *modules.ModuleServiceMap) *GetNewAddressService {
    return &GetNewAddressService{
        moduleServices,
    }
}

func (gnas *GetNewAddressService) GetNewAddressHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol      := vars["symbol"]
    addressType := vars["type"]

    SYMBOL := strings.ToUpper(symbol)

    logger.InfoLog(" - GetNewAddressHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "receiver")
    if err != nil {
        logger.ErrorLog(" - GetNewAddressHandler util.GetRpcConfigByType(SYMBOL, \"receiver\") err: "+err.Error())
        return
    }

    rpcRes, err := (*gnas.moduleServices)[SYMBOL].GetNewAddress(rpcConfig, addressType)
    if err != nil { 
        logger.ErrorLog(" - GetNewAddressHandler (*gnas.moduleServices)[SYMBOL].GetNewAddress(rpcConfig, addressType) addressType: "+addressType+", Error: "+err.Error())
        return
    }

    RES := GetNewAddressRes{
        RpcConfig: RpcConfigResDetail{ 
            RpcConfigId         : rpcConfig.Id,
            Symbol              : SYMBOL,
            Name                : rpcConfig.Name,
            Host                : rpcConfig.Host,
            Type                : rpcConfig.Type,
            NodeVersion         : rpcConfig.NodeVersion,
            NodeLastUpdated     : rpcConfig.NodeLastUpdated,
        },
        Address: rpcRes.Address,
    }

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetNewAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}
