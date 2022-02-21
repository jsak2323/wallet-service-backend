package coldbalance

import "context"

type Repository interface {
	Create(CreateColdBalance) (int, error)
	GetAll(ctx context.Context, page, limit int) ([]ColdBalance, error)
	GetByCurrencyId(ctx context.Context, currencyId int) ([]ColdBalance, error)
	GetByName(name string) (ColdBalance, error)
	GetByFireblocksName(ctx context.Context, name string) (ColdBalance, error)
	Update(ColdBalance) error
	UpdateBalance(ctx context.Context, id int, balance string) error
	ToggleActive(ctx context.Context, Id int, active bool) error
}
