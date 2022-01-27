package withdraw

type Withdraw struct {
	Id            int64
	CurrencyId    int
	AddressTo     string
	Tx            string
	Memo          string
	Amount        string
	Confirmations int
	BlockchainFee string
	MarketPrice   string
	LogIndex      string
	SuccessTime   string
	LastUpdated   string
}
