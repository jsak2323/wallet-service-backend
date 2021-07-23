package cold

import (
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
)

type ColdWalletService struct {
	cbRepo  cb.Repository
}

func NewColdWalletService(cbRepo cb.Repository) *ColdWalletService {
	return &ColdWalletService{cbRepo: cbRepo}
}