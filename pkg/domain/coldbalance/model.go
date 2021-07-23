package coldbalance

type ColdBalance struct {
	Id 			int 	`json:"id"`
	CurrencyId 	int 	`json:"currency_id"`
	Name 		string 	`json:"name"`
	Type 		string  `json:"type"`
	Balance 	float64 `json:"balance"`
	LastUpdated string	`json:"last_updated"`
}