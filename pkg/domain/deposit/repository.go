package deposit

import "context"

type Repository interface {
	CreateOrUpdate(Deposit) (int, error)
	Get(ctx context.Context, page, limit int, filters []map[string]interface{}) ([]Deposit, error)
	GetById(int) (Deposit, error)
	Update(Deposit) error
}
