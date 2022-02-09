package mysql

import (
	"context"
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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

func (r *rolePermissionRepository) Create(ctx context.Context, roleId, permissionId int) (err error) {
	query := "INSERT INTO " + rolePermissionTable + " (role_id, permission_id) VALUES(?, ?)"
	err = r.db.QueryRowContext(ctx, query, roleId, permissionId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rolePermissionRepository) GetByRole(roleId int) (rps []domain.RolePermission, err error) {
	rps, err = r.queryRows("SELECT role_id, permission_id FROM "+rolePermissionTable+" WHERE role_id = ?", roleId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *rolePermissionRepository) GetByPermission(permissionId int) (rps []domain.RolePermission, err error) {
	rps, err = r.queryRows("SELECT role_id, permission_id FROM "+rolePermissionTable+" WHERE permission_id = ?", permissionId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *rolePermissionRepository) queryRows(query string, param int) (rps []domain.RolePermission, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.RolePermission{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.RolePermission

		if err = rows.Scan(
			&rp.RoleId,
			&rp.PermissionId,
		); err != nil {
			return []domain.RolePermission{}, errs.AddTrace(err)
		}

		rps = append(rps, rp)
	}

	return rps, nil
}

func (r *rolePermissionRepository) DeleteByRoleId(ctx context.Context, roleId int) (err error) {
	query := "DELETE FROM " + rolePermissionTable + " WHERE role_id = ?"
	err = r.db.QueryRowContext(ctx, query, roleId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rolePermissionRepository) DeleteByPermissionId(ctx context.Context, permissionId int) (err error) {
	query := "DELETE FROM " + rolePermissionTable + " WHERE permission_id = ?"
	err = r.db.QueryRowContext(ctx, query, permissionId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rolePermissionRepository) Delete(roleId, permissionId int) (err error) {
	query := "DELETE FROM " + rolePermissionTable + " WHERE role_id = ? and permission_id = ?"
	err = r.db.QueryRow(query, roleId, permissionId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
