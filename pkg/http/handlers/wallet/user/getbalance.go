package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}

func (s *UserWalletService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, symbol string) {
	var (
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField))
		}
	}()

	if symbol == "" {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURRRPC))

		for _, SYMBOL := range config.SYMBOLS {
			go func(_SYMBOL string) {
				defer wg.Done()

				(*RES)[symbol], err = s.getUserBalanceRes(symbol)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetUserBalanceRes)
					return
				}
			}(SYMBOL)
		}

		wg.Wait()
	} else {
		(*RES)[symbol], err = s.getUserBalanceRes(symbol)
		if err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetUserBalanceRes)
			return
		}
	}
}

func (s *UserWalletService) getUserBalanceRes(symbol string) (TotalUserBalanceRes, error) {
	tcb, err := s.userBalanceRepo.GetTotalCoinBalance(symbol)
	if err != nil {
		return TotalUserBalanceRes{}, errs.AddTrace(err)
	}

	currencyConfigs, err := config.GetCurrencyBySymbol(symbol)
	if err != nil {
		return TotalUserBalanceRes{}, errs.AddTrace(err)
	}

	return TotalUserBalanceRes{
		TokenTypes: currencyConfigs,
		Balance:    tcb,
	}, nil
}
