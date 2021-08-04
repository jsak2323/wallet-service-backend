package cold

import (
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
)

const errInternalServer = "Internal server error"

type ColdWalletService struct {
	cbRepo  cb.Repository
}

func NewColdWalletService(cbRepo cb.Repository) *ColdWalletService {
	return &ColdWalletService{cbRepo: cbRepo}
}