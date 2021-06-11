package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
)

const permissionTable = "permissions"

type permissionRepository struct {
	db *sql.DB
}

func NewMysqlPermissionRepository(db *sql.DB) domain.Repository {
	return &permissionRepository{
		db,
	}
}

func (r *permissionRepository) Create(name string) (id int, err error) {
	query := "INSERT INTO " + permissionTable + " (name) VALUES(?)"

	err = r.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *permissionRepository) GetByName(name string) (permission domain.Permission, err error) {
	query := "SELECT id, name FROM " + permissionTable + " where name = ? limit 1"

	if err = r.db.QueryRow(query, name).Scan(
		&permission.Id,
		&permission.Name,
	); err != nil {
		return domain.Permission{}, err
	}

	return permission, nil
}
