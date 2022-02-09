package walletrpc

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/gorilla/mux"

// 	"github.com/btcid/wallet-services-backend-go/cmd/config"
// 	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
// 	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
// 	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
// 	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
// 	"github.com/btcid/wallet-services-backend-go/pkg/modules"
// )

// type GetBalanceHandlerResponseMap map[string][]GetBalanceRes

// type GetBalanceService struct {
// 	moduleServices *modules.ModuleServiceMap
// }

// func NewGetBalanceService(moduleServices *modules.ModuleServiceMap) *GetBalanceService {
// 	return &GetBalanceService{
// 		moduleServices,
// 	}
// }

// func (gbcs *GetBalanceService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	symbol := vars["symbol"]
// 	tokenType := vars["token_type"]
// 	isGetAll := symbol == ""
// 	ctx := req.Context()

// 	RES := make(GetBalanceHandlerResponseMap)

// 	if isGetAll {
// 		logger.InfoLog(" - GetBalanceHandler For all symbols, Requesting ...", req)
// 	} else {
// 		logger.InfoLog(" - GetBalanceHandler For symbol: "+strings.ToUpper(symbol)+", Requesting ...", req)
// 	}

// 	gbcs.InvokeGetBalance(ctx, &RES, symbol, tokenType)

// 	// handle success response
// 	resJson, _ := json.Marshal(RES)
// 	logger.InfoLog(" - GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(RES)
// }

// func (gbcs *GetBalanceService) InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, symbol, tokenType string) {
// 	var (
// 		errField       *errs.Error = nil
// 		rpcConfigCount             = 0
// 		resChannel                 = make(chan GetBalanceRes)
// 	)

// 	defer func() {
// 		if errField != nil {
// 			logger.ErrorLog(errs.Logged(errField), ctx)
// 		}
// 	}()

// 	for _, curr := range config.CURRRPC {
// 		SYMBOL := strings.ToUpper(curr.Config.Symbol)
// 		TOKENTYPE := strings.ToUpper(curr.Config.Symbol)

// 		// if symbol is defined, only get for that symbol
// 		if symbol != "" && strings.ToUpper(symbol) != SYMBOL && strings.ToUpper(tokenType) != TOKENTYPE {
// 			continue
// 		}

// 		for _, rpcConfig := range curr.RpcConfigs {
// 			rpcConfigCount++
// 			_RES := GetBalanceRes{
// 				RpcConfig: RpcConfigResDetail{
// 					RpcConfigId:          rpcConfig.Id,
// 					Name:                 rpcConfig.Name,
// 					Host:                 rpcConfig.Host,
// 					Type:                 rpcConfig.Type,
// 					NodeVersion:          rpcConfig.NodeVersion,
// 					NodeLastUpdated:      rpcConfig.NodeLastUpdated,
// 					IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
// 				},
// 			}

// 			// execute concurrent rpc calls
// 			go func(currencyConfig cc.CurrencyConfig, rpcConfig rc.RpcConfig) {
// 				module, err := gbcs.moduleServices.GetModule(currencyConfig.Id)
// 				if err != nil {
// 					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
// 					return
// 				}

// 				rpcRes, err := module.GetBalance(rpcConfig)
// 				if err != nil {
// 					_RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetBalance)

// 				} else {
// 					logger.Log(" -- InvokeGetBalance Symbol: "+SYMBOL+", RpcConfigId: "+strconv.Itoa(rpcConfig.Id)+", Host: "+rpcConfig.Host+". Balance: "+rpcRes.Balance, ctx)
// 					_RES.Balance = rpcRes.Balance
// 					_RES.Error = errs.AddTrace(errors.New(rpcRes.Error))
// 				}

// 				resChannel <- _RES

// 			}(curr.Config, rpcConfig)
// 		}
// 	}

// 	i := 0
// 	for res := range resChannel {
// 		i++
// 		(*RES)[res.RpcConfig.Symbol] = append((*RES)[res.RpcConfig.Symbol], res)
// 		if i >= rpcConfigCount {
// 			close(resChannel)
// 		}
// 	}
// }
