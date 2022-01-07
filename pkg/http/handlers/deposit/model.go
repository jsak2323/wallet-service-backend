package deposit

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type ListRes struct {
	Deposits []domain.Deposit `json:"deposits"`
	Error    *errs.Error      `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
