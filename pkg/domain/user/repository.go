package user

type Repository interface {
	Create(u User) (int, error)
	GetByUsername(username string) (User, error)
}
