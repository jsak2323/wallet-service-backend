package deposit

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
)

type ListRes struct {
	Deposits []domain.Deposit `json:"deposits"`
	Error    string           `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
