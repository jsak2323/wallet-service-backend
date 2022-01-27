package withdraw

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
)

type ListRes struct {
	Withdraws []domain.Withdraw `json:"withdraws"`
	Error    string           `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
