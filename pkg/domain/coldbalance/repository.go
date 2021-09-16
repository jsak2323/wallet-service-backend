package coldbalance

type Repository interface {
	GetAll(page, limit int) ([]ColdBalance, error)
	GetByCurrencyId(currencyId int) ([]ColdBalance, error)
	GetByName(name string) (ColdBalance, error)
	Update(ColdBalance) error
	UpdateBalance(id int, balance string) error
	ToggleActive(Id int, active bool) error
}