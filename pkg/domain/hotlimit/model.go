package hotlimit

const TopHardType = "top_hard"
const TopSoftType = "top_soft"
const TargetType = "target"
const BottomSoftType = "bottom_soft"
const BottomHardType = "bottom_hard"

type HotLimit struct {
	Id 			int 	`json:"id"`
	CurrencyId 	int 	`json:"currency_id"`
	Type 		string 	`json:"type"`
	Amount 		string 	`json:"amount"`
}