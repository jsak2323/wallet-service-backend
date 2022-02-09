package currencyconfig

import "context"

type Repository interface {
	Create(context.Context, CurrencyConfig) error
	GetAll(context.Context) ([]CurrencyConfig, error)
	GetBySymbol(symbol string) (*CurrencyConfig, error)
	Update(context.Context, UpdateCurrencyConfig) error
	ToggleActive(ctx context.Context, Id int, active bool) error
}
