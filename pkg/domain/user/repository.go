package user

type Repository interface {
	Create(User) (int, error)
	Update(User) error
	GetByUsername(username string) (User, error)
	GetEmailsByRole(role string) ([]string, error)
	GetAll(page, limit int) ([]User, error)
	ToggleActive(Id int, active bool) error
}
