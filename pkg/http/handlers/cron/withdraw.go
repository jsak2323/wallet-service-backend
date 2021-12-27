package cron

import (
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

	ltRES := make(h.ListWithdrawsHandlerResponseMap)

	logger.Log(" - Withdraw Update - Getting node list transactions ...")
	s.listWithdrawsService.InvokeListWithdraws(&ltRES, "", "", 1000)
	logger.Log(" - Withdraw Update - Getting node list transactions done. Fetched " + strconv.Itoa(len(ltRES)) + " results.")

	for symbol, mapTokenType := range ltRES {
		for tokenType, resRpcConfigs := range mapTokenType {
			var symbolTokenTypeLogStr = fmt.Sprintf("symbol: %s token_type: %s", symbol, tokenType)
			var txSuccessCount, txFailCount int
			var currency cc.CurrencyConfig
			var err error

			logger.Log(fmt.Sprintf("- Withdraw Update - Saving - %s...", symbolTokenTypeLogStr))

			if currency, err = config.GetCurrencyBySymbolTokenType(symbol, tokenType); err != nil {
				logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s err: %s", symbolTokenTypeLogStr, err.Error()))
				continue
			}

			for _, resRpcConfig := range resRpcConfigs {
				if resRpcConfig.Error != "" {
					logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s err: %s", symbolTokenTypeLogStr, resRpcConfig.Error))
					continue
				}
				
				logger.Log(fmt.Sprintf("- Withdraw Update - Saving - %s rpc_config_id: %d", symbolTokenTypeLogStr, resRpcConfig.RpcConfig.RpcConfigId))

				for _, tx := range resRpcConfig.Withdraws {
					amountRaw, err := util.CoinToRaw(tx.Amount, 8)
					if err != nil {
						logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s tx: %s value: %s util.CoinToRaw err: %s", symbolTokenTypeLogStr, tx.Tx, tx.Amount, err.Error()))
						continue
					}

					marketPrice, err := s.marketService.ConvertCoinToIdr(tx.Amount, symbol)
					if err != nil {
						logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s tx: %s value: %s marketService.ConvertCoinToIdr err: %s", symbolTokenTypeLogStr, tx.Tx, tx.Amount, err.Error()))
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

						logger.ErrorLog(fmt.Sprintf(" - Withdraw Update - %s tx: %s depositRepo.CreateOrUpdate err: %s", symbolTokenTypeLogStr, tx.Tx, err.Error()))
						continue
					}

					txSuccessCount++
				}
			}

			logger.Log(fmt.Sprintf("- Withdraw Update - Finished Saving - %s success: %d fail: %d", symbolTokenTypeLogStr, txSuccessCount, txFailCount))
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println(" - Withdraw Update - Time elapsed: " + fmt.Sprintf("%f", elapsedTime.Minutes()) + " minutes.")
}
