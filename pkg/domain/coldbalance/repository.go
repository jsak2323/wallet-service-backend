package coldbalance

type Repository interface {
	GetAll(page, limit int) ([]ColdBalance, error)
	GetByCurrencyId(currencyId int) ([]ColdBalance, error)
	GetByName(name string) (ColdBalance, error)
	GetDepositAddress(currencyId int, coldType string) (string, error)
	UpdateBalance(id int, balance string) error
}