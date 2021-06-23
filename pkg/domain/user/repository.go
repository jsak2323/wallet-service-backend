package user

type Repository interface {
	Create(User) (int, error)
	Update(User) error
	GetByUsername(username string) (User, error)
	GetAll(page, limit int) ([]User, error)
	ToggleActive(userId int, active bool) error
}
