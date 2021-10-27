package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
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
    tokenType  := vars["token_type"]
    address := vars["address"]

    SYMBOL := strings.ToUpper(symbol)
    logger.InfoLog(" - AddressTypeHandler For symbol: "+SYMBOL+", Requesting ...", req) 

    currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
    if err != nil {
        logger.ErrorLog(" - AddressTypeHandler config.GetCurrencyBySymbol("+symbol+","+tokenType+")+err: "+err.Error())
        RES.Error = err.Error()
        return
    }
    
    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(currencyConfig.Id, "receiver")
    if err != nil {
        logger.ErrorLog(" - AddressTypeHandler util.GetRpcConfigByType(SYMBOL, \"receiver\") err: "+err.Error())
        RES.Error = err.Error()
        return
    }
    RES.RpcConfig = RpcConfigResDetail{
        RpcConfigId             : rpcConfig.Id,
        Name                    : rpcConfig.Name,
        Host                    : rpcConfig.Host,
        Type                    : rpcConfig.Type,
        NodeVersion             : rpcConfig.NodeVersion,
        NodeLastUpdated         : rpcConfig.NodeLastUpdated,
        IsHealthCheckEnabled    : rpcConfig.IsHealthCheckEnabled,
    }

    module, err := ats.moduleServices.GetModule(currencyConfig.Id)
    if err != nil {
        logger.ErrorLog(" - AddressTypeHandler stas.moduleServices.GetModule err: "+err.Error())
        RES.Error = err.Error()
        return
    }

    // execute rpc call
    rpcRes, err := module.AddressType(rpcConfig, address)
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


