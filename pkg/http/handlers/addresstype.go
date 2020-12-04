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
    // define response object
    RES := AddressTypeRes{}

    // define response handler
    handleResponse := func() {
        resStatus := http.StatusOK
        if RES.Error != "" {
            resStatus = http.StatusInternalServerError
        }
        w.WriteHeader(resStatus)
        json.NewEncoder(w).Encode(RES)
    }
    defer handleResponse()

    // define request params
    vars := mux.Vars(req)
    symbol  := vars["symbol"]
    address := vars["address"]

    SYMBOL := strings.ToUpper(symbol)
    logger.InfoLog(" - AddressTypeHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "receiver")
    if err != nil {
        logger.ErrorLog(" - AddressTypeHandler util.GetRpcConfigByType(SYMBOL, \"receiver\") err: "+err.Error())
        RES.Error = err.Error()
        return
    }
    RES.RpcConfig = RpcConfigResDetail{
        RpcConfigId             : rpcConfig.Id,
        Symbol                  : SYMBOL,
        Name                    : rpcConfig.Name,
        Host                    : rpcConfig.Host,
        Type                    : rpcConfig.Type,
        NodeVersion             : rpcConfig.NodeVersion,
        NodeLastUpdated         : rpcConfig.NodeLastUpdated,
        IsHealthCheckEnabled    : rpcConfig.IsHealthCheckEnabled,
    }

    // execute rpc call
    rpcRes, err := (*ats.moduleServices)[SYMBOL].AddressType(rpcConfig, address)
    if err != nil { 
        logger.ErrorLog(" - AddressTypeHandler (*ats.moduleServices)[SYMBOL].AddressType(rpcConfig, addressType) address: "+address+", Error: "+err.Error())
        RES.Error = err.Error()
        return
    }

    // handle success response
    RES.AddressType = rpcRes.AddressType
    RES.Error       = rpcRes.Error
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - AddressTypeHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}


