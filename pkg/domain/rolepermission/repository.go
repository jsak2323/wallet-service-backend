package rolepermission

import "context"

type Repository interface {
	Create(ctx context.Context, roleId, permissionId int) error
	GetByRole(roleId int) ([]RolePermission, error)
	GetByPermission(permissionId int) ([]RolePermission, error)
	DeleteByRoleId(ctx context.Context, roleId int) error
	DeleteByPermissionId(ctx context.Context, permissionId int) error
	Delete(roleId, permissionId int) error
}
