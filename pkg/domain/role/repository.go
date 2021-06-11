package role

type Repository interface {
	Create(name string) (id int, err error)
	GetByID(id int) (Role, error)
	GetByName(name string) (Role, error)
}
