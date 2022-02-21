package wallet

import (
	"context"
	"sync"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

func (s *walletService) GetBalance(ctx context.Context, currConfig config.CurrencyRpcConfig) GetBalanceRes {
	var wg sync.WaitGroup
	var res GetBalanceRes = GetBalanceRes{CurrencyConfig: currConfig.Config}

	wg.Add(5)
	go func() { defer wg.Done(); s.SetColdBalanceDetails(ctx, &res) }()
	go func() { defer wg.Done(); s.SetHotBalanceDetails(ctx, currConfig.RpcConfigs, &res) }()
	go func() { defer wg.Done(); s.SetUserBalanceDetails(ctx, &res) }()
	go func() { defer wg.Done(); s.SetPendingWithdraw(ctx, &res) }()
	go func() { defer wg.Done(); s.SetHotLimits(ctx, &res) }()
	wg.Wait()

	s.SetPercent(ctx, &res)

	return res
}
