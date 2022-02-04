package userrole

import "context"

type Repository interface {
	Create(ctx context.Context, userId, roleId int) error
	GetByUser(ctx context.Context, userId int) ([]UserRole, error)
	GetByRole(ctx context.Context, roleId int) ([]UserRole, error)
	DeleteByUserId(ctx context.Context, userId int) error
	DeleteByRoleId(ctx context.Context, roleId int) error
	Delete(ctx context.Context, userId, roleId int) error
}
