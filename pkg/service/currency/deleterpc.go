package currency

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) DeleteCurrencyRpc(ctx context.Context, currencyId int, rpcId int) (err error) {
	if err = s.crRepo.Delete(ctx, currencyId, rpcId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedDeleteCurrencyRPC)
		return err
	}
	return nil
}
