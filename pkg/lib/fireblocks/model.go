package fireblocks

type GetVaultAccountAssetReq struct {
	VaultAccountId 	string 	`json:"vault_account_id"`
	AssetId 		string 	`json:"asset_id"`
}

type GetVaultAccountAssetRes struct {
	Total string `json:"total"`
	Error 	string `json:"message"`
}

type CreateTransactionReq struct {
	AssetId 	string 				`json:"asset_id"`
	Source 		TransactionAccount 	`json:"source"`
	Amount 		string 				`json:"amount"`
	Destination TransactionAccount 	`json:"destination"`
}

type TransactionAccount struct {
	Type 	string `json:"type"`
	Id 		string `json:"id"`
}

type CreateTransactionRes struct {
	Id		string `json:"id"`
	Status	string `json:"status"`
	Error 	string `json:"message"`
}