package user

import (
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
)

type UserWalletService struct {
	userBalanceRepo ub.Repository
}

func NewUserWalletService(ubRepo ub.Repository) *UserWalletService {
	return &UserWalletService{ubRepo}
}