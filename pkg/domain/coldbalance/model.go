package coldbalance

const FbColdType = "fb_cold"
const FbWarmType = "fb_warm"
const ColdType = "cold"
const Cold2Type = "cold_2"
const StakingType = "staking"

type ColdBalance struct {
	Id 					int
	CurrencyId 			int
	Name 				string
	Type 				string
	FireblocksName 		string
	Balance 			string
	Address				string
	LastUpdated 		string
	Active 				bool
}