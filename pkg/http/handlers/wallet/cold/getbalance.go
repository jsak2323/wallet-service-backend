package cold

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceResponse struct {
	Currency cc.CurrencyConfig
	Balances []cb.ColdBalance
}

type GetBalanceHandlerResponseMap map[string]GetBalanceResponse

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
				
				(*RES)[_SYMBOL] = GetBalanceResponse{
					Currency: currConfig,
					Balances: s.getBalance(_SYMBOL, currConfig),
				}
			}(curr.Config, SYMBOL)
		}

		wg.Wait()
    } else {
		(*RES)[symbol] = GetBalanceResponse{
			Currency: config.CURR[symbol].Config,
			Balances: s.getBalance(symbol, config.CURR[symbol].Config),
		}
	}
}

func (s *ColdWalletService) getBalance(symbol string, currency cc.CurrencyConfig) (coldBalances []cb.ColdBalance) {
	if currency.FireblocksName != "" {
		if res, err := fireblocks.GetVaultAccountAsset(fireblocks.GetVaultAccountAssetReq{
			VaultAccountId: config.CONF.FireblocksColdVaultId,
			AssetId: currency.FireblocksName,
		}); err != nil {
			logger.ErrorLog("- cold.getBalance fireblocks.GetVaultAccountAsset("+currency.FireblocksName+") error: "+err.Error())
		} else {
			coldBalance := cb.ColdBalance{
				Name: currency.Symbol + " Cold",
				Type: "cold",
				CurrencyId: currency.Id,
			}
			if coldBalance.Balance, err = strconv.ParseFloat(res.Balance, 64); err != nil {
				logger.ErrorLog("- cold.getBalance strconv.ParseFloat("+currency.FireblocksName+") error: "+err.Error())
			}
			coldBalances = append(coldBalances, coldBalance)
		}
	}

	if cbs, err := s.cbRepo.GetByCurrencyId(currency.Id); err != nil {
		logger.ErrorLog("- cold.getBalance s.cbRepo.GetByCurrencyId("+strconv.Itoa(currency.Id)+") error: "+err.Error())
	} else if len(cbs) > 0 {
		coldBalances = append(coldBalances, cbs...)
	}

	return coldBalances
}