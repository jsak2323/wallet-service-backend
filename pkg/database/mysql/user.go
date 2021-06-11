package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/user"
)

const userTable = "users"

type userRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.Repository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(u domain.User) (id int, err error) {
	query := "INSERT INTO " + userTable + " (username, name, password, ip_address) VALUES(?, ?, ?, ?)"

	err = r.db.QueryRow(query, u.Username, u.Name, u.Password, u.IPAddress).Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepository) GetByUsername(username string) (user domain.User, err error) {
	query := "SELECT * FROM " + userTable + " where username = ? limit 1"

	rows, err := r.db.Query(query, username)
	defer rows.Close()
	if err != nil {
		return domain.User{}, err
	}

	for rows.Next() {
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Password,
			&user.IPAddress,
		)
		if err != nil {
			return domain.User{}, err
		}
	}

	return user, nil
}
