package role

import "context"

type Repository interface {
	Create(ctx context.Context, name string) (id int, err error)
	Update(context.Context, Role) error
	GetAll(page, limit int) ([]Role, error)
	GetByName(name string) (Role, error)
	GetById(id int) (Role, error)
	GetByUserId(userId int) ([]Role, error)
	GetNamesByUserId(ctx context.Context, userId int) ([]string, error)
	Delete(ctx context.Context, roleId int) error
}
