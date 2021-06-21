package mysql

import (
	"strconv"
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/role"
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

func (r *roleRepository) Create(name string) (id int, err error) {
	query := "INSERT INTO " + roleTable + " (name) VALUES(?);"

	if err = r.db.QueryRow(query, name).Err(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *roleRepository) Update(role domain.Role) (err error) {
	query := "UPDATE " + roleTable + " SET name = ? WHERE id = ?"

	err = r.db.QueryRow(query, role.Name, role.Id).Err()
	if err != nil {
		return err
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
		return []domain.Role{}, err
	}
	defer rows.Close()

	for rows.Next() {
		role := domain.Role{}

		if err = rows.Scan(&role.Id,&role.Name); err != nil {
			return []domain.Role{}, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) GetByName(name string) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where name = ? limit 1"

	if err = r.db.QueryRow(query, name).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) GetById(id int) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where id = ? limit 1"

	if err = r.db.QueryRow(query, id).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) GetByUserId(userId int) (roles []domain.Role, err error) {
	query := "SELECT r.id, r.name FROM " + roleTable + " as r"
	query = query + " JOIN user_role ur ON ur.role_id = r.id"
	query = query + " WHERE ur.user_id = ?"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []domain.Role{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var role domain.Role

		if err = rows.Scan(&role.Id, &role.Name); err != nil {
			return []domain.Role{}, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) GetNamesByUserId(userId int) (roles []string, err error) {
	query := "SELECT r.name FROM " + roleTable + " as r"
	query = query + " JOIN user_role ur ON ur.role_id = r.id"
	query = query + " WHERE ur.user_id = ?"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		if err = rows.Scan(&name); err != nil {
			return []string{}, err
		}

		roles = append(roles, name)
	}

	return roles, nil
}

func(r *roleRepository) Delete(roleId int) error {
	query := "DELETE FROM " + roleTable + " WHERE id = ?"

	return r.db.QueryRow(query, roleId).Err()
}