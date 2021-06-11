package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
)

const rolePermissionTable = "role_permission"

type rolePermissionRepository struct {
	db *sql.DB
}

func NewMysqlRolePermissionRepository(db *sql.DB) domain.Repository {
	return &rolePermissionRepository{
		db,
	}
}

func (r *rolePermissionRepository) Create(roleId, permissionId int) (err error) {
	query := "INSERT INTO " + rolePermissionTable + " (role_id, permission_id) VALUES(?, ?)"

	return r.db.QueryRow(query, roleId, permissionId).Err()
}

func (r *rolePermissionRepository) GetByRole(roleId int) (rps []domain.RolePermission, err error) {
	return r.queryRows("SELECT role_id, permission_id FROM "+rolePermissionTable+" WHERE role_id = ?", roleId)
}

func (r *rolePermissionRepository) GetByPermission(permissionId int) (rps []domain.RolePermission, err error) {
	return r.queryRows("SELECT role_id, permission_id FROM "+rolePermissionTable+" WHERE permission_id = ?", permissionId)
}

func (r *rolePermissionRepository) queryRows(query string, param int) (rps []domain.RolePermission, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.RolePermission{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.RolePermission

		if err = rows.Scan(
			&rp.RoleId,
			&rp.PermissionId,
		); err != nil {
			return []domain.RolePermission{}, err
		}

		rps = append(rps, rp)
	}

	return rps, nil
}
