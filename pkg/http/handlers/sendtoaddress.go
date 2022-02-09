package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	config "github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type SendToAddressRequest struct {
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Memo      string `json:"memo"`
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
	ctx := req.Context()

	// define response handler
	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	// define request params
	sendToAddressRequest := SendToAddressRequest{}
	err := DecodeAndLogPostRequest(req, &sendToAddressRequest)
	if err != nil {
		logger.ErrorLog(" - SendToAddressHandler util.DecodeAndLogPostRequest(req, &sendToAddressRequest) err: "+err.Error(), ctx)
		return
	}

	symbol := sendToAddressRequest.Symbol
	tokenType := sendToAddressRequest.TokenType
	amountInDecimal := sendToAddressRequest.Amount
	address := sendToAddressRequest.Address
	memo := sendToAddressRequest.Memo

	SYMBOL := strings.ToUpper(symbol)
	TOKENTYPE := strings.ToUpper(tokenType)
	logger.InfoLog(" - SendToAddressHandler Sending "+amountInDecimal+" "+SYMBOL+" "+TOKENTYPE+", Requesting ...", req)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return
	}

	// define rpc config
	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "sender")
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return
	}
	RES.RpcConfig = RpcConfigResDetail{
		RpcConfigId:          rpcConfig.Id,
		Symbol:               SYMBOL,
		TokenType:            TOKENTYPE,
		Name:                 rpcConfig.Name,
		Host:                 rpcConfig.Host,
		Type:                 rpcConfig.Type,
		NodeVersion:          rpcConfig.NodeVersion,
		NodeLastUpdated:      rpcConfig.NodeLastUpdated,
		IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
	}

	module, err := stas.moduleServices.GetModule(currencyConfig.Id)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
		return
	}

	// execute rpc call
	rpcRes, err := module.SendToAddress(ctx, rpcConfig, amountInDecimal, address, memo)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedSendToAddress)
		return
	}

	// handle success response
	RES.TxHash = rpcRes.TxHash
	RES.Error = errs.AddTrace(errors.New(rpcRes.Error))
	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - SendToAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}
