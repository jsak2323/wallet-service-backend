package cron

import (
	"fmt"
	"strconv"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	d "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type DepositService struct {
    moduleServices  *modules.ModuleServiceMap
    listTransactionsService  *h.ListTransactionsService
    depositRepo d.Repository
}

func NewDepositService(
    moduleServices  *modules.ModuleServiceMap,
    listTransactionsService  *h.ListTransactionsService,
    depositRepo d.Repository,
) *DepositService {
    return &DepositService{
        moduleServices,
		listTransactionsService,
        depositRepo,
    }
}

func (s *DepositService) Update() {
	startTime := time.Now()

	ltRES := make(h.ListTransactionsHandlerResponseMap)

    logger.Log(" - Deposit Update - Getting node list transactions ...")
	s.listTransactionsService.InvokeListTransactions(&ltRES, "", "", 1000)
    logger.Log(" - Deposit Update - Getting node list transactions done. Fetched "+strconv.Itoa(len(ltRES))+" results." )

	for symbol, mapTokenType := range ltRES { 
		for tokenType, resRpcConfigs := range mapTokenType { 
			var symbolTokenTypeLogStr = fmt.Sprintf("symbol: %s token_type: %s", symbol, tokenType)
			var txSuccessCount, txFailCount int
			var currency cc.CurrencyConfig
			var err error

			logger.Log(fmt.Sprintf("- Deposit Update - Saving - %s...", symbolTokenTypeLogStr))			
			
			if currency, err = config.GetCurrencyBySymbolTokenType(symbol, tokenType); err != nil {
				logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s err: %s", symbolTokenTypeLogStr, err.Error()));
				continue
			}
				
			for _, resRpcConfig := range resRpcConfigs {
				logger.Log(fmt.Sprintf("- Deposit Update - Saving - %s rpc_config_id: %d", symbolTokenTypeLogStr, resRpcConfig.RpcConfig.RpcConfigId))

				for _, tx := range resRpcConfig.Transactions {
					amountRaw, err := util.CoinToRaw(tx.Amount, 8)
					if err != nil {
						logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s tx: %s value: %s util.CoinToRaw err: %s", symbolTokenTypeLogStr, tx.Tx, tx.Amount, err.Error()));
						continue
					}
					
					if _, err = s.depositRepo.CreateOrUpdate(d.Deposit{
						CurrencyId: currency.Id,
						AddressTo: tx.AddressTo,
						Tx: tx.Tx,
						Amount: amountRaw,
						Memo: tx.Memo,
						SuccessTime: tx.SuccessTime,
					}); err != nil {
						txFailCount++
						
						logger.ErrorLog(fmt.Sprintf(" - Deposit Update - %s tx: %s depositRepo.CreateOrUpdate err: %s", symbolTokenTypeLogStr, tx.Tx, err.Error()));
						continue
					}

					txSuccessCount++
				}
			}

			logger.Log(fmt.Sprintf("- Deposit Update - Finished Saving - %s success: %d fail: %d", symbolTokenTypeLogStr, txSuccessCount, txFailCount))			
		}
	}

    elapsedTime := time.Since(startTime)
    fmt.Println(" - Deposit Update - Time elapsed: "+fmt.Sprintf("%f", elapsedTime.Minutes())+ " minutes.")
}