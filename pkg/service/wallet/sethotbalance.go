package wallet

import (
	"context"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	modulesm "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (s *walletService) SetHotBalanceDetails(ctx context.Context, rpcConfigs []rc.RpcConfig, res *GetBalanceRes) {
	var (
		symbol   string      = res.CurrencyConfig.Symbol
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	for _, rpcConfig := range rpcConfigs {
		var hotBalanceDetail BalanceDetail = BalanceDetail{Name: rpcConfig.Name, Type: rpcConfig.Type}
		var rpcRes *modulesm.GetBalanceRpcRes

		module, err := s.moduleServices.GetModule(res.CurrencyConfig.Id)
		if err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
			continue
		}

		if rpcRes, err = module.GetBalance(ctx, rpcConfig); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
			continue
		}

		hotBalanceDetail.Coin = rpcRes.Balance

		if hotBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(hotBalanceDetail.Coin, symbol); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		if res.TotalHotCoin, err = util.AddCurrency(res.TotalHotCoin, hotBalanceDetail.Coin); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		if res.TotalHotIdr, err = util.AddCurrency(res.TotalHotIdr, hotBalanceDetail.Idr); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		res.HotBalances = append(res.HotBalances, hotBalanceDetail)
	}
}
