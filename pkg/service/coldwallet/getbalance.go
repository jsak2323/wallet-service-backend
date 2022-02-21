package coldwallet

import (
	"context"
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *coldWalletService) GetBalance(ctx context.Context, currencyConfigId int) (coldBalances []cb.ColdBalance) {

	var (
		currency             = config.CURRRPC[currencyConfigId].Config
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if cbs, err := s.cbRepo.GetByCurrencyId(ctx, currency.Id); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyByID)
	} else if len(cbs) > 0 {
		for i := range cbs {
			if cbs[i].Type == cb.FbWarmType || cbs[i].Type == cb.FbColdType {
				vaultAccountId, err := s.FireblocksVaultAccountId(cbs[i].Type)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.FailedFireblocksVaultAccountId)
				}

				if res, err := fireblocks.GetVaultAccountAsset(fireblocks.GetVaultAccountAssetReq{
					VaultAccountId: vaultAccountId,
					AssetId:        cbs[i].FireblocksName,
				}); err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.FailedGetVaultAccountAsset)
				} else {
					cbs[i].Balance = res.Total
				}
			} else {
				// non fireblocks balance are stored in raw in db
				cbs[i].Balance = util.RawToCoin(cbs[i].Balance, 8)
			}
		}
		coldBalances = append(coldBalances, cbs...)
	}

	return coldBalances
}

func (s *coldWalletService) FireblocksVaultAccountId(cbType string) (string, error) {
	switch cbType {
	case domain.FbColdType:
		return config.CONF.FireblocksColdVaultId, nil
	case domain.FbWarmType:
		return config.CONF.FireblocksWarmVaultId, nil
	}

	return "", errs.AddTrace(errors.New("invalid fireblocks type: " + cbType))
}
