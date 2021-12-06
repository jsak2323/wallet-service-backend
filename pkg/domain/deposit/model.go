package deposit

type Deposit struct {
	Id          int64
	CurrencyId  int
	AddressTo   string
	Tx          string
	Memo        string
	Amount      string
	SuccessTime string
	LastUpdated string
}
