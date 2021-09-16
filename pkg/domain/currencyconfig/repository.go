package currencyconfig

type CurrencyConfigRepository interface {
    Create(CurrencyConfig) error
    GetAll() ([]CurrencyConfig, error)
    GetBySymbol(symbol string) (*CurrencyConfig, error)
    Update(CurrencyConfig) error
	ToggleActive(Id int, active bool) error
}


