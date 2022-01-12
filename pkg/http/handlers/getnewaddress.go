package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	config "github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
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
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	// define request params
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	tokenType := vars["token_type"]
	addressType := vars["type"]

	SYMBOL := strings.ToUpper(symbol)
	logger.InfoLog(" - GetNewAddressHandler For symbol: "+SYMBOL+", Requesting ...", req)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return
	}

	// define rpc config
	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "receiver")
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return
	}
	RES.RpcConfig = RpcConfigResDetail{
		RpcConfigId:          rpcConfig.Id,
		Symbol:               SYMBOL,
		Name:                 rpcConfig.Name,
		Host:                 rpcConfig.Host,
		Type:                 rpcConfig.Type,
		NodeVersion:          rpcConfig.NodeVersion,
		NodeLastUpdated:      rpcConfig.NodeLastUpdated,
		IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
	}

	module, err := gnas.moduleServices.GetModule(currencyConfig.Id)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
		return
	}

	// execute rpc call
	rpcRes, err := module.GetNewAddress(rpcConfig, addressType)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetNewAddress)
		return
	}

	// handle success response
	RES.Address = rpcRes.Address
	RES.Error = errs.AddTrace(errors.New(rpcRes.Error))
	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - GetNewAddressHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
}
