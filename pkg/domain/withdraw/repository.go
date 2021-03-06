package withdraw

import "context"

type Repository interface {
	CreateOrUpdate(Withdraw) (int, error)
	Get(ctx context.Context, page, limit int, filters []map[string]interface{}) ([]Withdraw, error)
	GetById(int) (Withdraw, error)
	Update(Withdraw) error
}
