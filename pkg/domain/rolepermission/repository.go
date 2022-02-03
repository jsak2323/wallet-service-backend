package rolepermission

import "context"

type Repository interface {
	Create(roleId, permissionId int) error
	GetByRole(roleId int) ([]RolePermission, error)
	GetByPermission(permissionId int) ([]RolePermission, error)
	DeleteByRoleId(roleId int) error
	DeleteByPermissionId(ctx context.Context, permissionId int) error
	Delete(roleId, permissionId int) error
}
