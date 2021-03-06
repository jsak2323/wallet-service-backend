package cold

import errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"

type SendToHotReq struct {
	FireblocksName string `json:"fireblocks_name"`
	FireblocksType string `json:"fireblocks_type"`
	Amount         string `json:"amount"`
	Memo           string `json:"memo"`
}

type UpdateBalanceReq struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
