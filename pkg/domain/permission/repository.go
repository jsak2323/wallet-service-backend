package permission

type Repository interface {
	Create(name string) (id int, err error)
	GetByName(name string) (Permission, error)
}
