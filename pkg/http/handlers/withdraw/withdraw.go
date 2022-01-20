package withdraw

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
)

const errInternalServer = "Internal server error"

type WithdrawService struct {
	wRepo domain.Repository
}

func NewWithdrawService(
	wRepo domain.Repository,
) *WithdrawService {
	return &WithdrawService{
		wRepo: wRepo,
	}
}
