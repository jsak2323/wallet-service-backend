package coldwallet

import (
	"context"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *coldWalletService) ListColdWallet(ctx context.Context, page int, limit int) (resp []cb.ColdBalance, err error) {
	if resp, err = s.cbRepo.GetAll(ctx, page, limit); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllColdBalance)
		return resp, err
	}
	return resp, nil
}
