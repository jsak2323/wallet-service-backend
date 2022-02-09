package currency

import (
	"context"

	currencyHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *currencyService) CreateCurrencyRpc(ctx context.Context, cRreq currencyHandler.CurrencyRpcReq) (err error) {
	if err = s.validator.Validate(cRreq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.crRepo.Create(ctx, cRreq.CurrencyId, cRreq.RpcId); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateCurrencyRPC)
		return err
	}
	return nil
}
