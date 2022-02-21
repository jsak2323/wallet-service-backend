package coldwallet

import (
	"context"
	"errors"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *coldWalletService) CreateColdWallet(ctx context.Context, createReq domain.CreateColdBalance) (err error) {
	if err = s.validator.Validate(createReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	if err = s.validateCreateReq(createReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return err
	}

	createReq.Balance = util.CoinToRaw(createReq.Balance, 8)

	if _, err = s.cbRepo.Create(createReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateColdBalance)
		return err
	}
	return nil
}

func (s *coldWalletService) validateCreateReq(createReq domain.CreateColdBalance) error {

	if s.isFireblocksCold(createReq.Type) && createReq.FireblocksName == "" {
		return errs.AddTrace(errors.New("Fireblocks Name"))
	}

	return nil
}

func (s *coldWalletService) isFireblocksCold(cbType string) bool {
	if cbType == domain.FbColdType || cbType == domain.FbWarmType {
		return true
	}

	return false
}
