package rolepermission

type Repository interface {
	Create(roleId, permissionId int) error
	GetByRole(roleId int) ([]RolePermission, error)
	GetByPermission(permissionId int) ([]RolePermission, error)
}
