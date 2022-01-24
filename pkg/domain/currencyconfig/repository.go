package currencyconfig

type Repository interface {
	Create(CurrencyConfig) error
	GetAll() ([]CurrencyConfig, error)
	GetBySymbol(symbol string) (*CurrencyConfig, error)
	Update(UpdateCurrencyConfig) error
	ToggleActive(Id int, active bool) error
}
