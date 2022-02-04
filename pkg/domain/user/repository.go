package user

import "context"

type Repository interface {
	Create(context.Context, User) (int, error)
	Update(context.Context, User) error
	GetByUsername(ctx context.Context, username string) (User, error)
	GetEmailsByRole(role string) ([]string, error)
	GetAll(page, limit int) ([]User, error)
	ToggleActive(ctx context.Context, Id int, active bool) error
}
