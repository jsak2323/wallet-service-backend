package permission

import "context"

type Repository interface {
	Create(ctx context.Context, name string) (id int, err error)
	Update(context.Context, Permission) error
	GetAll(ctx context.Context, page, limit int) ([]Permission, error)
	GetByName(name string) (Permission, error)
	GetByRoleId(roleId int) ([]Permission, error)
	GetNamesByUserId(userId int) ([]string, error)
	GetNamesByRoleId(roleId int) ([]string, error)
	Delete(ctx context.Context, permissionId int) error
}
