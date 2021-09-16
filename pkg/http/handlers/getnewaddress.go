package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
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
    // define response object
    RES := GetNewAddressRes{}

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
    symbol      := vars["symbol"]
    addressType := vars["type"]

    SYMBOL := strings.ToUpper(symbol)
    logger.InfoLog(" - GetNewAddressHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "receiver")
    if err != nil {
        logger.ErrorLog(" - GetNewAddressHandler util.GetRpcConfigByType(SYMBOL, \"receiver\") err: "+err.Error())
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

    module, ok := (*gnas.moduleServices)[SYMBOL]
    if !ok {
        logger.ErrorLog(" - GetNewAddressHandler module not implemented symbol: "+SYMBOL)
    }
    
    // execute rpc call
    rpcRes, err := module.GetNewAddress(rpcConfig, addressType)
    if err != nil { 
        logger.ErrorLog(" - GetNewAddressHandler (*gnas.moduleServices)[SYMBOL].GetNewAddress(rpcConfig, addressType) addressType: "+addressType+", Error: "+err.Error())
        RES.Error = err.Error()
        return
    }

    // handle success response
    RES.Address = rpcRes.Address
    RES.Error   = rpcRes.Error
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - GetNewAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}


