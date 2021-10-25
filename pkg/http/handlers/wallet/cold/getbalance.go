package cold

import (
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) GetBalance(symbol string) (coldBalances []cb.ColdBalance) {
	currency := config.CURR[symbol].Config

	if cbs, err := s.cbRepo.GetByCurrencyId(currency.Id); err != nil {
		logger.ErrorLog("- cold.getBalance s.cbRepo.GetByCurrencyId(" + strconv.Itoa(currency.Id) + ") error: " + err.Error())
	} else if len(cbs) > 0 {
		for i := range cbs {
			if cbs[i].Type == cb.FbWarmType || cbs[i].Type == cb.FbColdType {
				vaultAccountId, err := FireblocksVaultAccountId(cbs[i].Type)
				if err != nil {
					logger.ErrorLog(" - cold.getBalance FireblocksVaultAccountId err: " + err.Error())
				}

				if res, err := fireblocks.GetVaultAccountAsset(fireblocks.GetVaultAccountAssetReq{
					VaultAccountId: vaultAccountId,
					AssetId:        cbs[i].FireblocksName,
				}); err != nil {
					logger.ErrorLog("- cold.getBalance fireblocks.GetVaultAccountAsset(" + cbs[i].FireblocksName + ") error: " + err.Error())
				} else {
					cbs[i].Balance = res.Total
				}
			} else {
				// non fireblocks balance are stored in raw in db
				if cbs[i].Balance, err = util.RawToCoin(cbs[i].Balance, 8); err != nil {
					logger.ErrorLog("- cold.getBalance RawToCoin(" + strconv.Itoa(currency.Id) + "," + cbs[i].Balance + ") error: " + err.Error())
				}
			}
		}
		coldBalances = append(coldBalances, cbs...)
	}

	return coldBalances
}
