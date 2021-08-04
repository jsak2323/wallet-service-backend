package balance

import (
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
)

type BalanceService struct {
	walletService *hw.WalletService
}

func NewBalanceService(
	walletService *hw.WalletService,
) *BalanceService {
	return &BalanceService{
		walletService,
	}
}