package deposit

type Deposit struct {
	Id            int64
	CurrencyId    int
	AddressTo     string
	Tx            string
	Memo          string
	Amount        string
	Confirmations int
	LogIndex      int
	SuccessTime   string
	LastUpdated   string
}
