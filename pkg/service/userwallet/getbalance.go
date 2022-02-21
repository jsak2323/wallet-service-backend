package userwallet

import (
	"context"
	"sync"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	handlers "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type GetBalanceHandlerResponseMap map[string]handlers.TotalUserBalanceRes

func (s *userWalletService) InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, symbol string) {
	var (
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
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

func (s *userWalletService) getUserBalanceRes(symbol string) (resp handlers.TotalUserBalanceRes, err error) {
	tcb, err := s.userBalanceRepo.GetTotalCoinBalance(symbol)
	if err != nil {
		return resp, errs.AddTrace(err)
	}

	currencyConfigs, err := config.GetCurrencyBySymbol(symbol)
	if err != nil {
		return resp, errs.AddTrace(err)
	}

	resp = handlers.TotalUserBalanceRes{
		TokenTypes: currencyConfigs,
		Balance:    tcb,
	}

	return resp, nil

}
