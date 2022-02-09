package currency

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	currencyHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) ListCurrency(ctx context.Context) (RES currencyHandler.ListRes, err error) {
	if len(config.CURRRPC) > 0 {
		for _, curr := range config.CURRRPC {
			RES.CurrencyConfigs = append(RES.CurrencyConfigs, curr.Config)
		}
	} else {
		if RES.CurrencyConfigs, err = s.ccRepo.GetAll(ctx); err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllCurrencyConfig)
			return RES, err
		}
	}

	for i, currency := range RES.CurrencyConfigs {
		RES.CurrencyConfigs[i].RpcConfigs, err = s.rcRepo.GetByCurrencyId(ctx, currency.Id)
		if err != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCConfigByCurrencyID)
			return RES, err
		}
	}
	return RES, nil
}
