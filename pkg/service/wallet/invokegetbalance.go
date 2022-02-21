package wallet

import (
	"context"
	"sync"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

func (s *walletService) InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, currencyId int) {
	if currencyId == 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURRRPC))

		for _, curr := range config.CURRRPC {
			go func(currencyConfiguration config.CurrencyRpcConfig) {
				defer wg.Done()

				(*RES)[currencyId] = s.GetBalance(ctx, currencyConfiguration)
			}(curr)
		}

		wg.Wait()
	} else {
		(*RES)[currencyId] = s.GetBalance(ctx, config.CURRRPC[currencyId])
	}
}
