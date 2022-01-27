package withdraw

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type ListRes struct {
	Withdraws []domain.Withdraw `json:"withdraws"`
	Error     *errs.Error       `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
