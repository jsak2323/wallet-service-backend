package mysql

import (
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

func (r *roleRepository) GetByName(name string) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where name = ? limit 1"

	if err = r.db.QueryRow(query, name).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) GetByID(id int) (role domain.Role, err error) {
	query := "SELECT id, name FROM " + roleTable + " where id = ? limit 1"

	if err = r.db.QueryRow(query, id).Scan(&role.Id, &role.Name); err != nil {
		return domain.Role{}, err
	}

	return role, nil
}
