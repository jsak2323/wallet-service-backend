package cold

import (
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
)

const errInternalServer = "Internal server error"

type ColdWalletService struct {
	cbRepo domain.Repository
}

func NewColdWalletService(cbRepo domain.Repository) *ColdWalletService {
	return &ColdWalletService{cbRepo: cbRepo}
}

func FireblocksVaultAccountId(cbType string) (string, error) {
	switch cbType {
	case domain.FbColdType:
		return config.CONF.FireblocksColdVaultId, nil
	case domain.FbWarmType:
		return config.CONF.FireblocksWarmVaultId, nil
	}

	return "", errors.New("invalid fireblocks type: " + cbType)
}

func isFireblocksCold(cbType string) bool {
	if cbType == domain.FbColdType || cbType == domain.FbWarmType {
		return true
	}

	return false
}
