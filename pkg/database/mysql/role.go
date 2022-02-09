package mysql

import (
	"context"
	"database/sql"
	"strconv"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const roleTable = "roles"

type roleRepository struct {
	db *sql.DB
}

func NewMysqlRoleRepository(db *sql.DB) domain.Repository {
	return &roleRepository{
		db,
	}
}

func (r *roleRepository) Create(ctx context.Context, name string) (id int, err error) {
	query := "INSERT INTO " + roleTable + " (name) VALUES(?);"

	if err = r.db.QueryRowContext(ctx, query, name).Err(); err != nil {
		return 0, errs.AddTrace(err)
	}

	return id, nil
}

func (r *roleRepository) Update(ctx context.Context, role domain.Role) (err error) {
	query := "UPDATE " + roleTable + " SET name = ? WHERE id = ?"

	err = r.db.QueryRowContext(ctx, query, role.Name, role.Id).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *roleRepository) GetAll(page, limit int) (roles []domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable

	if limit <= 0 {
		limit = defaultLimit
	}

	if page > 0 {
		query = query + " offset " + strconv.Itoa(page) + " limit " + strconv.Itoa(limit)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return []domain.Role{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		role := domain.Role{}

		if err = rows.Scan(&role.Id, &role.Name); err != nil {
			return []domain.Role{}, errs.AddTrace(err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) GetByName(name string) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where name = ? limit 1"

	if err = r.db.QueryRow(query, name).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, errs.AddTrace(err)
	}

	return role, nil
}

func (r *roleRepository) GetById(id int) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where id = ? limit 1"

	if err = r.db.QueryRow(query, id).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, errs.AddTrace(err)
	}

	return role, nil
}

func (r *roleRepository) GetByUserId(userId int) (roles []domain.Role, err error) {
	query := "SELECT r.id, r.name FROM " + roleTable + " as r"
	query = query + " JOIN user_role ur ON ur.role_id = r.id"
	query = query + " WHERE ur.user_id = ?"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []domain.Role{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var role domain.Role

		if err = rows.Scan(&role.Id, &role.Name); err != nil {
			return []domain.Role{}, errs.AddTrace(err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) GetNamesByUserId(ctx context.Context, userId int) (roles []string, err error) {
	query := "SELECT r.name FROM " + roleTable + " as r"
	query = query + " JOIN user_role ur ON ur.role_id = r.id"
	query = query + " WHERE ur.user_id = ?"

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return []string{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		if err = rows.Scan(&name); err != nil {
			return []string{}, errs.AddTrace(err)
		}

		roles = append(roles, name)
	}

	return roles, nil
}

func (r *roleRepository) Delete(ctx context.Context, roleId int) error {
	query := "DELETE FROM " + roleTable + " WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, roleId).Err()
	if err != nil {
		errs.AddTrace(err)

	}

	return errs.AddTrace(err)
}
