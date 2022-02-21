package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	d "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/service/walletrpc"
)

func (s *cronService) UpdateDeposit() {
	startTime := time.Now()
	ctx := context.Background()

	ltRES := make(walletrpc.ListTransactionsHandlerResponseMap)

	logger.Log(" - Deposit Update - Getting node list transactions ...", ctx)
	s.walletRpcService.InvokeListTransactions(ctx, &ltRES, "", "", 1000)

	logger.Log(" - Deposit Update - Getting node list transactions done. Fetched "+strconv.Itoa(len(ltRES))+" results.", ctx)

	for symbol, mapTokenType := range ltRES {
		for tokenType, resRpcConfigs := range mapTokenType {
			var symbolTokenTypeLogStr = fmt.Sprintf("symbol: %s token_type: %s", symbol, tokenType)
			var txSuccessCount, txFailCount int
			var currency cc.CurrencyConfig
			var err error

			logger.Log(fmt.Sprintf("- Deposit Update - Saving - %s...", symbolTokenTypeLogStr), ctx)

			if currency, err = config.GetCurrencyBySymbolTokenType(symbol, tokenType); err != nil {
				logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s err: %s", symbolTokenTypeLogStr, err.Error()), ctx)
				continue
			}

			for _, resRpcConfig := range resRpcConfigs {
				if resRpcConfig.Error != nil {
					logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s err: %s", symbolTokenTypeLogStr, resRpcConfig.Error), ctx)
					continue
				}
				logger.Log(fmt.Sprintf("- Deposit Update - Saving - %s rpc_config_id: %d", symbolTokenTypeLogStr, resRpcConfig.RpcConfig.RpcConfigId), ctx)

				for _, tx := range resRpcConfig.Transactions {
					amountRaw := util.CoinToRaw(tx.Amount, 8)

					if _, err = s.depositRepo.CreateOrUpdate(d.Deposit{
						CurrencyId:  currency.Id,
						AddressTo:   tx.AddressTo,
						Tx:          tx.Tx,
						Amount:      amountRaw,
						Memo:        tx.Memo,
						SuccessTime: tx.SuccessTime,
					}); err != nil {
						txFailCount++

						logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s tx: %s depositRepo.CreateOrUpdate err: %s", symbolTokenTypeLogStr, tx.Tx, err.Error()), ctx)
						continue
					}

					txSuccessCount++
				}
			}

			logger.Log(fmt.Sprintf("- Deposit Update - Finished Saving - %s success: %d fail: %d", symbolTokenTypeLogStr, txSuccessCount, txFailCount), ctx)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println(" - Deposit Update - Time elapsed: " + fmt.Sprintf("%f", elapsedTime.Minutes()) + " minutes.")
}
