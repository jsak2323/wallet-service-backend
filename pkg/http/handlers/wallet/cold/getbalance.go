package cold

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceRes struct {
	Currency cc.CurrencyConfig
	Balances []cb.ColdBalance
}

type GetBalanceHandlerResponseMap map[string]GetBalanceRes

func (s *ColdWalletService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    symbol := strings.ToUpper(vars["symbol"])
    isGetAll := symbol == ""
	
	RES := make(GetBalanceHandlerResponseMap)

	if isGetAll {
        logger.InfoLog(" - cold.GetBalanceHandler For all symbols, Requesting ...", req)
	} else {
        logger.InfoLog(" - cold.GetBalanceHandler For symbol: "+symbol+", Requesting ...", req)
	}
    
	s.invokeGetBalance(&RES, symbol)

	resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - cold.GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (s *ColdWalletService) invokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol string) {
	if symbol == "" {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURR))

		for SYMBOL, curr := range config.CURR {
			go func(currConfig cc.CurrencyConfig, _SYMBOL string) {
				defer wg.Done()
				
				(*RES)[_SYMBOL] = GetBalanceRes{
					Currency: currConfig,
					Balances: s.GetBalance(_SYMBOL),
				}
			}(curr.Config, SYMBOL)
		}

		wg.Wait()
    } else {
		(*RES)[symbol] = GetBalanceRes{
			Currency: config.CURR[symbol].Config,
			Balances: s.GetBalance(symbol),
		}
	}
}

func (s *ColdWalletService) GetBalance(symbol string) (coldBalances []cb.ColdBalance) {
	currency := config.CURR[symbol].Config

	if cbs, err := s.cbRepo.GetByCurrencyId(currency.Id); err != nil {
		logger.ErrorLog("- cold.getBalance s.cbRepo.GetByCurrencyId("+strconv.Itoa(currency.Id)+") error: "+err.Error())
	} else if len(cbs) > 0 {
		logger.ErrorLog(cbs[0].Name + " " + cbs[0].Type + " " + cbs[0].FireblocksName)
		for i := range cbs {
			if cbs[i].Type == cb.FbWarmType || cbs[i].Type == cb.FbColdType {
				if res, err := fireblocks.GetVaultAccountAsset(fireblocks.GetVaultAccountAssetReq{
					VaultAccountId: FireblocksVaultAccountId(cbs[i].Type),
					AssetId: cbs[i].FireblocksName,
				}); err != nil {
					logger.ErrorLog("- cold.getBalance fireblocks.GetVaultAccountAsset("+cbs[i].FireblocksName+") error: "+err.Error())
				} else {
					cbs[i].Balance = res.Balance
				}
			} else {
				// non fireblocks balance are stored in raw in db
				if cbs[i].Balance, err = util.RawToCoin(cbs[i].Balance, 8); err != nil {
					logger.ErrorLog("- cold.getBalance RawToCoin("+strconv.Itoa(currency.Id)+","+cbs[i].Balance+") error: "+err.Error())
				}
			}
		}
		coldBalances = append(coldBalances, cbs...)
	}

	return coldBalances
}