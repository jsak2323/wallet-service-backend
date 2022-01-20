package withdraw

type Repository interface {
	CreateOrUpdate(Withdraw) (int, error)
	Get(page, limit int, filters []map[string]interface{}) ([]Withdraw, error)
	GetById(int) (Withdraw, error)
	Update(Withdraw) error
}
