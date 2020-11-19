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

type AddressTypeService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewAddressTypeService(moduleServices *modules.ModuleServiceMap) *AddressTypeService {
    return &AddressTypeService{
        moduleServices,
    }
}

func (ats *AddressTypeService) AddressTypeHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    symbol  := vars["symbol"]
    address := vars["address"]

    SYMBOL := strings.ToUpper(symbol)

    logger.InfoLog(" - AddressTypeHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "receiver")
    if err != nil {
        logger.ErrorLog(" - AddressTypeHandler util.GetRpcConfigByType(SYMBOL, \"receiver\") err: "+err.Error())
        return
    }

    rpcRes, err := (*ats.moduleServices)[SYMBOL].AddressType(rpcConfig, address)
    if err != nil { 
        logger.ErrorLog(" - AddressTypeHandler (*ats.moduleServices)[SYMBOL].AddressType(rpcConfig, addressType) address: "+address+", Error: "+err.Error())
        return
    }

    RES := AddressTypeRes{
        RpcConfig: RpcConfigResDetail{ 
            RpcConfigId         : rpcConfig.Id,
            Symbol              : SYMBOL,
            Name                : rpcConfig.Name,
            Host                : rpcConfig.Host,
            Type                : rpcConfig.Type,
            NodeVersion         : rpcConfig.NodeVersion,
            NodeLastUpdated     : rpcConfig.NodeLastUpdated,
        },
        AddressType : rpcRes.AddressType,
        Error       : rpcRes.Error,
    }

    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - AddressTypeHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}


