package mysql

import (
	"database/sql"
	"strconv"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/user"
)

const userTable = "users"
const userTableAlias = "u"
const defaultLimit = 10

type userRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.Repository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(u domain.User) (id int, err error) {
	query := "INSERT INTO " + userTable + " (username, name, email, password, ip_address) VALUES(?, ?, ?, ?, ?)"

	err = r.db.QueryRow(query, u.Username, u.Name, u.Email, u.Password, u.IPAddress).Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepository) Update(u domain.User) (err error) {
	var params []interface{}
	var query  string

	query = "UPDATE " + userTable + " SET name = ?, username = ?, email = ?, ip_address = ?"
	params = append(params, u.Name, u.Username, u.Email, u.IPAddress)

	if u.Password != "" {
		query = query + ", password = ? "
		params = append(params, u.Password)
	}

	query = query + " WHERE id = ?"
	params = append(params, u.Id)

	if err = r.db.QueryRow(query, params...).Err(); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByUsername(username string) (user domain.User, err error) {
	query := "SELECT id, username, name, email, password, ip_address, active FROM " + userTable + " where username = ? limit 1"

	err = r.db.QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IPAddress,
		&user.Active,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r userRepository) GetAll(page, limit int) (users []domain.User, err error) {
	query := "SELECT id, username, name, email, password, ip_address, active FROM " + userTable
	
	if limit <= 0 {
		limit = defaultLimit
	}
	
	if page > 0 {
		query = query + " offset " + strconv.Itoa(page) + " limit " + strconv.Itoa(limit)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return []domain.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		user := domain.User{}

		if err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.IPAddress,
			&user.Active,
		); err != nil {
			return []domain.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func(r *userRepository) ToggleActive(userId int, active bool) error {
	query := "UPDATE " + userTable + " SET active = ? WHERE id = ?"

	return r.db.QueryRow(query, active, userId).Err()
}
