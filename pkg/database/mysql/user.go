package mysql

import (
	"context"
	"database/sql"
	"strconv"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const userTable = "users"
const defaultLimit = 10

type userRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.Repository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(ctx context.Context, u domain.User) (id int, err error) {
	query := "INSERT INTO " + userTable + " (username, name, email, password, ip_address) VALUES(?, ?, ?, ?, ?)"

	err = r.db.QueryRowContext(ctx, query, u.Username, u.Name, u.Email, u.Password, u.IPAddress).Err()
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	return id, nil
}

func (r *userRepository) Update(ctx context.Context, u domain.User) (err error) {
	var params []interface{}
	var query string

	query = "UPDATE " + userTable + " SET name = ?, username = ?, email = ?, ip_address = ?"
	params = append(params, u.Name, u.Username, u.Email, u.IPAddress)

	if u.Password != "" {
		query = query + ", password = ? "
		params = append(params, u.Password)
	}

	query = query + " WHERE id = ?"
	params = append(params, u.Id)

	if err = r.db.QueryRowContext(ctx, query, params...).Err(); err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (user domain.User, err error) {
	query := "SELECT id, username, name, email, password, ip_address, active FROM " + userTable + " where username = ? limit 1"

	err = r.db.QueryRowContext(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IPAddress,
		&user.Active,
	)
	if err != nil {
		return domain.User{}, errs.AddTrace(err)
	}

	return user, nil
}

func (r userRepository) GetEmailsByRole(role string) (emails []string, err error) {
	query := "SELECT email FROM " + userTable
	query += " JOIN " + userRoleTable + " on " + userRoleTable + ".user_id = " + userTable + ".id"
	query += " JOIN " + roleTable + " on " + roleTable + ".id = " + userRoleTable + ".role_id"
	query += " WHERE " + roleTable + ".name = ?"

	rows, err := r.db.Query(query, role)
	if err != nil {
		return []string{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		email := ""

		if err = rows.Scan(
			&email,
		); err != nil {
			return []string{}, errs.AddTrace(err)
		}

		emails = append(emails, email)
	}

	return emails, nil
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
		return []domain.User{}, errs.AddTrace(err)
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
			return []domain.User{}, errs.AddTrace(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) ToggleActive(ctx context.Context, userId int, active bool) error {
	query := "UPDATE " + userTable + " SET active = ? WHERE id = ?"

	if err := r.db.QueryRowContext(ctx, query, active, userId).Err(); err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
