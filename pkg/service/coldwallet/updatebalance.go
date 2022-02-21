package coldwallet

import (
	"context"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *coldWalletService) UpdateBalance(ctx context.Context, id int, balance string) (err error) {
	balanceRaw := util.CoinToRaw(balance, 8)

	if err = s.cbRepo.UpdateBalance(ctx, id, balanceRaw); err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
