package cold

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
)

func(s *ColdWalletService) GetDepositAddress(currencyId int) (string, error) {
	return s.cbRepo.GetDepositAddress(currencyId, domain.ColdType)
}