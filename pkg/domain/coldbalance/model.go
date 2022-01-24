package coldbalance

const FbColdType = "fb_cold"
const FbWarmType = "fb_warm"
const ColdType = "cold"
const Cold2Type = "cold_2"
const StakingType = "staking"

type ColdBalance struct {
	Id int `validate:"required"`
	CreateColdBalance
}

type CreateColdBalance struct {
	Id             int
	CurrencyId     int    `validate:"required"`
	Name           string `validate:"required"`
	Type           string `validate:"required"`
	FireblocksName string
	Balance        string
	Address        string `validate:"required"`
	LastUpdated    string
	Active         bool
}
