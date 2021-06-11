package userrole

type Repository interface {
	Create(userId, roleId int) error
	GetByUser(userId int) ([]UserRole, error)
	GetByRole(roleId int) ([]UserRole, error)
}
