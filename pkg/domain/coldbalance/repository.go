package coldbalance

import "context"

type Repository interface {
	Create(CreateColdBalance) (int, error)
	GetAll(page, limit int) ([]ColdBalance, error)
	GetByCurrencyId(currencyId int) ([]ColdBalance, error)
	GetByName(name string) (ColdBalance, error)
	GetByFireblocksName(ctx context.Context, name string) (ColdBalance, error)
	Update(ColdBalance) error
	UpdateBalance(id int, balance string) error
	ToggleActive(Id int, active bool) error
}
