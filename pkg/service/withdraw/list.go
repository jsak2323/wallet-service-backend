package withdraw

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *withdrawService) List(ctx context.Context, page int, limit int, filters []map[string]interface{}) (resp []withdraw.Withdraw, err error) {

	if resp, err = s.withdrawRepo.Get(ctx, page, limit, filters); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetListWithdraws)
		return resp, err
	}

	return resp, nil
}
