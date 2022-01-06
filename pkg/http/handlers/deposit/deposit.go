package deposit

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
)

const errInternalServer = "Internal server error"

type DepositService struct {
	dRepo domain.Repository
}

func NewDepositService(
	dRepo domain.Repository,
) *DepositService {
	return &DepositService{
		dRepo: dRepo,
	}
}
