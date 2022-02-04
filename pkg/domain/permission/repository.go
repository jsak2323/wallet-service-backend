package permission

import "context"

type Repository interface {
	Create(ctx context.Context, name string) (id int, err error)
	Update(context.Context, Permission) error
	GetAll(ctx context.Context, page, limit int) ([]Permission, error)
	GetByName(ctx context.Context, name string) (Permission, error)
	GetByRoleId(ctx context.Context, roleId int) ([]Permission, error)
	GetNamesByUserId(ctx context.Context, userId int) ([]string, error)
	GetNamesByRoleId(ctx context.Context, roleId int) ([]string, error)
	Delete(ctx context.Context, permissionId int) error
}
