package permission

type Repository interface {
	Create(name string) (id int, err error)
	Update(Permission) error
	GetAll(page, limit int) ([]Permission, error)
	GetByName(name string) (Permission, error)
	GetByRoleId(roleId int) ([]Permission, error)
	GetNamesByUserId(userId int) ([]string, error)
	GetNamesByRoleId(roleId int) ([]string, error)
	Delete(permissionId int) (error)
}
