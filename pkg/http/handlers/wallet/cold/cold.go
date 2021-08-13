package cold

import (
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

const errInternalServer = "Internal server error"

type ColdWalletService struct {
	cbRepo  cb.Repository
}

func NewColdWalletService(cbRepo cb.Repository) *ColdWalletService {
	return &ColdWalletService{cbRepo: cbRepo}
}

func FireblocksVaultAccountId(cbType string) string {
	switch cbType {
		case cb.FbColdType: return config.CONF.FireblocksColdVaultId
		case cb.FbWarmType: return config.CONF.FireblocksWarmVaultId
	}

	return config.CONF.FireblocksColdVaultId
}