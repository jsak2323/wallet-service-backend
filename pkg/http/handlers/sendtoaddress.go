package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type SendToAddressRequest struct {
    Symbol      string `json:"symbol"` 
    Amount      string `json:"amount"` 
    Address     string `json:"address"` 
    Memo        string `json:"memo"` 
}

type SendToAddressService struct {
    moduleServices *modules.ModuleServiceMap
}

func NewSendToAddressService(moduleServices *modules.ModuleServiceMap) *SendToAddressService {
    return &SendToAddressService{
        moduleServices,
    }
}

func (stas *SendToAddressService) SendToAddressHandler(w http.ResponseWriter, req *http.Request) { 
    // define response object
    RES := SendToAddressRes{}

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
    sendToAddressRequest := SendToAddressRequest{}
    err := DecodeAndLogPostRequest(req, &sendToAddressRequest)
    if err != nil {
        logger.ErrorLog(" - SendToAddressHandler util.DecodeAndLogPostRequest(req, &sendToAddressRequest) err: "+err.Error())
        return
    }
    symbol          := sendToAddressRequest.Symbol
    amountInDecimal := sendToAddressRequest.Amount
    address         := sendToAddressRequest.Address
    memo            := sendToAddressRequest.Memo

    SYMBOL := strings.ToUpper(symbol)
    logger.InfoLog(" - SendToAddressHandler Sending "+amountInDecimal+" "+SYMBOL+", Requesting ...", req) 

    // define rpc config
    rpcConfig, err := util.GetRpcConfigByType(SYMBOL, "sender")
    if err != nil {
        logger.ErrorLog(" - SendToAddressHandler util.GetRpcConfigByType(SYMBOL, \"sender\") err: "+err.Error())
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
    rpcRes, err := (*stas.moduleServices)[SYMBOL].SendToAddress(rpcConfig, amountInDecimal, address, memo)
    if err != nil { 
        logger.ErrorLog(" - SendToAddressHandler (*stas.moduleServices)[strings.ToUpper(symbol)].SendToAddress(rpcConfig, address, amountInDecimal) address:"+address+", amount: "+amountInDecimal+", Error: "+err.Error())
        RES.Error = err.Error()
        return
    }

    // handle success response
    RES.TxHash = rpcRes.TxHash
    RES.Error  = rpcRes.Error
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - SendToAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}


