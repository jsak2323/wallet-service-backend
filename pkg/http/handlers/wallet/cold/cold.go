package cold

import (
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type ColdWalletService struct {
	cbRepo    domain.Repository
	validator util.CustomValidator
}

func NewColdWalletService(cbRepo domain.Repository, validator util.CustomValidator) *ColdWalletService {
	return &ColdWalletService{cbRepo: cbRepo, validator: validator}
}

func FireblocksVaultAccountId(cbType string) (string, error) {
	switch cbType {
	case domain.FbColdType:
		return config.CONF.FireblocksColdVaultId, nil
	case domain.FbWarmType:
		return config.CONF.FireblocksWarmVaultId, nil
	}

	return "", errs.AddTrace(errors.New("invalid fireblocks type: " + cbType))
}

func isFireblocksCold(cbType string) bool {
	if cbType == domain.FbColdType || cbType == domain.FbWarmType {
		return true
	}

	return false
}
