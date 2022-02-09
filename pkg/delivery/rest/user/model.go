package user

import errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
