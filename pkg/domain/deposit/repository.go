package deposit

type Repository interface {
	Create(Deposit) (int, error)
	Get(page, limit int, filters []map[string]interface{}) ([]Deposit, error)
	GetById(int) (Deposit, error)
	Update(Deposit) error
}
