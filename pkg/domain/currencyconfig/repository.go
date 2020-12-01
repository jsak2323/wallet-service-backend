package currencyconfig

type CurrencyConfigRepository interface {
    GetAll() ([]CurrencyConfig, error)
    GetBySymbol(symbol string) (*CurrencyConfig, error)
}


