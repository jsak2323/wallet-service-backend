package cold

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *ColdWalletService) SettlementWallet(currencyId int) (result domain.ColdBalance, err error) {
	cbs, err := s.cbRepo.GetByCurrencyId(currencyId)
	if err != nil {
		return domain.ColdBalance{}, errs.AddTrace(err)
	}

	for _, cb := range cbs {
		// prioritise on fb_warm type, break loop when found
		if cb.Type == domain.FbWarmType {
			result = cb
			break
		}

		// take cold type in case there's no fb_warm type
		if cb.Type == domain.ColdType {
			result = cb
		}
	}

	return result, nil
}
