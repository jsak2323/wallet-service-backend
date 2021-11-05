package user

import (
	"net/http"
	"encoding/json"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
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
	var err error

	if symbol == "" {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURRRPC))

		for _, SYMBOL := range config.SYMBOLS {
			go func(_SYMBOL string) {
				defer wg.Done()

				(*RES)[symbol], err = s.getUserBalanceRes(symbol)
				if err != nil {
					logger.ErrorLog("- userwallet.GetBalance s.userBalanceRepo.GetTotalCoinBalance("+symbol+") error: "+err.Error())
					return
				}
			}(SYMBOL)
		}

		wg.Wait()
	} else {
		(*RES)[symbol], err = s.getUserBalanceRes(symbol)
		if err != nil {
			logger.ErrorLog("- userwallet.GetBalance s.userBalanceRepo.GetTotalCoinBalance("+symbol+") error: "+err.Error())
			return
		}
	}
}

func (s *UserWalletService) getUserBalanceRes(symbol string) (TotalUserBalanceRes, error) {
	tcb, err := s.userBalanceRepo.GetTotalCoinBalance(symbol)
		if err != nil {
			return TotalUserBalanceRes{}, err
		} 

		currencyConfigs, err := config.GetCurrencyBySymbol(symbol)
		if err != nil {
			return TotalUserBalanceRes{}, err
		}
		
		return TotalUserBalanceRes{
			TokenTypes: currencyConfigs,
			Balance: tcb,
		}, nil
}