package walletrpc

// import (
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"github.com/gorilla/mux"

// 	"github.com/btcid/wallet-services-backend-go/cmd/config"
// 	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
// 	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
// 	"github.com/btcid/wallet-services-backend-go/pkg/modules"
// )

// type AddressTypeService struct {
// 	moduleServices *modules.ModuleServiceMap
// 	toogle         bool
// }

// func NewAddressTypeService(moduleServices *modules.ModuleServiceMap, toogle bool) *AddressTypeService {
// 	return &AddressTypeService{
// 		moduleServices,
// 		toogle,
// 	}
// }

// func (ats *AddressTypeService) AddressTypeHandler(w http.ResponseWriter, req *http.Request) {
// 	// define response object
// 	RES := AddressTypeRes{}
// 	ctx := req.Context()

// 	// define response handler
// 	handleResponse := func() {
// 		resStatus := http.StatusOK
// 		if RES.Error != nil {
// 			resStatus = http.StatusInternalServerError
// 			logger.ErrorLog(errs.Logged(RES.Error), ctx)
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(resStatus)
// 		json.NewEncoder(w).Encode(RES)
// 	}
// 	defer handleResponse()

// 	// define request params
// 	vars := mux.Vars(req)
// 	symbol := vars["symbol"]
// 	tokenType := vars["token_type"]
// 	address := vars["address"]

// 	SYMBOL := strings.ToUpper(symbol)
// 	logger.InfoLog(" - AddressTypeHandler For symbol: "+SYMBOL+", Requesting ...", req)

// 	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
// 	if err != nil {
// 		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
// 		return
// 	}

// 	// define rpc config
// 	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "receiver")
// 	if err != nil {
// 		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
// 		return
// 	}
// 	RES.RpcConfig = RpcConfigResDetail{
// 		RpcConfigId:          rpcConfig.Id,
// 		Name:                 rpcConfig.Name,
// 		Host:                 rpcConfig.Host,
// 		Type:                 rpcConfig.Type,
// 		NodeVersion:          rpcConfig.NodeVersion,
// 		NodeLastUpdated:      rpcConfig.NodeLastUpdated,
// 		IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
// 	}

// 	module, err := ats.moduleServices.GetModule(currencyConfig.Id)
// 	if err != nil {
// 		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
// 		return
// 	}

// 	// execute rpc call
// 	rpcRes, err := module.AddressType(ctx, rpcConfig, address)
// 	if err != nil {
// 		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedAddressType)
// 		return
// 	}

// 	// handle success response
// 	RES.AddressType = rpcRes.AddressType
// 	RES.Error = errs.AddTrace(errors.New(rpcRes.Error))
// 	resJson, _ := json.Marshal(RES)
// 	logger.InfoLog(" - AddressTypeHandler Success. Symbol: "+SYMBOL+", Res: "+string(resJson), req)
// }
