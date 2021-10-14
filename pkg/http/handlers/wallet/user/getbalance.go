package user

import (
	"net/http"
	"encoding/json"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceHandlerResponseMap map[string]TotalUserBalanceRes

func (s *UserWalletService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    symbol := strings.ToUpper(vars["symbol"])
	isGetAll := symbol != ""

	RES := make(GetBalanceHandlerResponseMap)

	if isGetAll {
        logger.InfoLog(" - userwallet.GetBalanceHandler For all symbols, Requesting ...", req)
	} else {
        logger.InfoLog(" - userwallet.GetBalanceHandler For symbol: "+symbol+", Requesting ...", req) 
	}
	
	s.InvokeGetBalance(&RES, symbol)
	
	resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - userwallet.GetBalanceHandler Success. Symbol: "+symbol+", Res: "+string(resJson), req)
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (s *UserWalletService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol string) {
	if symbol == "" {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURR))

		for _, curr := range config.CURR {
			go func(_curr cc.CurrencyConfig) {
				defer wg.Done()

				if tcb, err := s.userBalanceRepo.GetTotalCoinBalance(_curr.Symbol); err != nil {
					logger.ErrorLog("- userwallet.GetBalance s.userBalanceRepo.GetTotalCoinBalance("+_curr.Symbol+") error: "+err.Error())
				} else {
					(*RES)[_curr.Symbol] = TotalUserBalanceRes{_curr, tcb}
				}
			}(curr.Config)
		}

		wg.Wait()
	} else {
		if tcb, err := s.userBalanceRepo.GetTotalCoinBalance(symbol); err != nil {
			logger.ErrorLog("- userwallet.GetBalance s.userBalanceRepo.GetTotalCoinBalance("+symbol+") error: "+err.Error())
		} else {
			(*RES)[symbol] = TotalUserBalanceRes{config.CURR[symbol].Config, tcb}
		}
	}
}