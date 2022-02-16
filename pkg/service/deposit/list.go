package deposit

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *depositService) List(ctx context.Context, page int, limit int, filters []map[string]interface{}) (resps []deposit.Deposit, err error) {
	if resps, err = s.depositRepo.Get(ctx, page, limit, filters); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetDeposit)
		return resps, err
	}

	return resps, nil
}
