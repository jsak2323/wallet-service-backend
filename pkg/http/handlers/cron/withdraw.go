package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	d "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type WithdrawService struct {
	moduleServices       *modules.ModuleServiceMap
	listWithdrawsService *h.ListWithdrawsService
	marketService        *h.MarketService
	depositRepo          d.Repository
}

func NewWithdrawService(
	moduleServices *modules.ModuleServiceMap,
	listWithdrawsService *h.ListWithdrawsService,
	marketService *h.MarketService,
	depositRepo d.Repository,
) *WithdrawService {
	return &WithdrawService{
		moduleServices,
		listWithdrawsService,
		marketService,
		depositRepo,
	}
}

func (s *WithdrawService) Update() {
	startTime := time.Now()
	ctx := context.Background()

	ltRES := make(h.ListWithdrawsHandlerResponseMap)

	logger.Log(" - Withdraw Update - Getting node list transactions ...", ctx)
	s.listWithdrawsService.InvokeListWithdraws(ctx, &ltRES, "", "", 1000)
	logger.Log(" - Withdraw Update - Getting node list transactions done. Fetched "+strconv.Itoa(len(ltRES))+" results.", ctx)

	for symbol, mapTokenType := range ltRES {
		for tokenType, resRpcConfigs := range mapTokenType {
			var symbolTokenTypeLogStr = fmt.Sprintf("symbol: %s token_type: %s", symbol, tokenType)
			var txSuccessCount, txFailCount int
			var currency cc.CurrencyConfig
			var err error

			logger.Log(fmt.Sprintf("- Withdraw Update - Saving - %s...", symbolTokenTypeLogStr), ctx)

			if currency, err = config.GetCurrencyBySymbolTokenType(symbol, tokenType); err != nil {
				logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s err: %s", symbolTokenTypeLogStr, err.Error()), ctx)
				continue
			}

			for _, resRpcConfig := range resRpcConfigs {
				if resRpcConfig.Error != "" {
					logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s err: %s", symbolTokenTypeLogStr, resRpcConfig.Error), ctx)
					continue
				}

				logger.Log(fmt.Sprintf("- Withdraw Update - Saving - %s rpc_config_id: %d", symbolTokenTypeLogStr, resRpcConfig.RpcConfig.RpcConfigId), ctx)

				for _, tx := range resRpcConfig.Withdraws {
					amountRaw := util.CoinToRaw(tx.Amount, 8)

					marketPrice, err := s.marketService.ConvertCoinToIdr(tx.Amount, symbol)
					if err != nil {
						logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s tx: %s value: %s marketService.ConvertCoinToIdr err: %s", symbolTokenTypeLogStr, tx.Tx, tx.Amount, err.Error()), ctx)
						continue
					}

					if _, err = s.depositRepo.CreateOrUpdate(d.Withdraw{
						CurrencyId:    currency.Id,
						AddressTo:     tx.AddressTo,
						Tx:            tx.Tx,
						Amount:        amountRaw,
						Confirmations: tx.Confirmations,
						BlockchainFee: tx.BlockchainFee,
						MarketPrice:   marketPrice,
						Memo:          tx.Memo,
						SuccessTime:   tx.SuccessTime,
					}); err != nil {
						txFailCount++

						logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s tx: %s depositRepo.CreateOrUpdate err: %s", symbolTokenTypeLogStr, tx.Tx, err.Error()), ctx)
						continue
					}

					txSuccessCount++
				}
			}

			logger.Log(fmt.Sprintf("- Withdraw Update - Finished Saving - %s success: %d fail: %d", symbolTokenTypeLogStr, txSuccessCount, txFailCount), ctx)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println(" - Withdraw Update - Time elapsed: " + fmt.Sprintf("%f", elapsedTime.Minutes()) + " minutes.")
}
