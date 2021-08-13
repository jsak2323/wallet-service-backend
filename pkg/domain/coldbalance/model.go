package coldbalance

const FbColdType = "fb_cold"
const FbWarmType = "fb_warm"
const ColdType = "cold"
const Cold2Type = "cold_2"
const StakingType = "staking"

type ColdBalance struct {
	Id 					int 	`json:"id"`
	CurrencyId 			int 	`json:"currency_id"`
	Name 				string 	`json:"name"`
	Type 				string  `json:"type"`
	FireblocksName 		string  `json:"fireblocks_name"`
	Balance 			string  `json:"balance"`
	Address				string 	`json:"address"`
	LastUpdated 		string	`json:"last_updated"`
}