package fireblocks

type GetVaultAccountAssetReq struct {
	VaultAccountId 	int 	`json:"vault_account_id"`
	AssetId 		string 	`json:"asset_id"`
}

type GetVaultAccountAssetRes struct {
	Balance string `json:"balance"`
	Error 	string `json:"message"`
}

type CreateTransactionReq struct {
	AssetId 	string 				`json:"asset_id"`
	Source 		TransactionAccount 	`json:"source"`
	Amount 		float64 			`json:"amount"`
	Destination TransactionAccount 	`json:"destination"`
}

type TransactionAccount struct {
	Type 	string `json:"type"`
	Id 		string `json:"id"`
}

type CreateTransactionRes struct {
	Status	string `json:"status"`
	Error 	string `json:"message"`
}