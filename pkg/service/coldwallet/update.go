package coldwallet

import (
	"context"
	"errors"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *coldWalletService) UpdateColdWallet(ctx context.Context, updateReq domain.ColdBalance) (err error) {
	if err = s.validator.Validate(updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if s.isFireblocksCold(updateReq.Type) && updateReq.FireblocksName == "" {
		return errs.AddTrace(errors.New("Fireblocks Name"))
	}

	updateReq.Balance = util.CoinToRaw(updateReq.Balance, 8)

	if err = s.cbRepo.Update(updateReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateColdBalance)
		return err
	}
	return nil
}
