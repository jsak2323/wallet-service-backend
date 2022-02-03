package mysql

import (
	"context"
	"database/sql"
	"strconv"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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

func (r *permissionRepository) Create(ctx context.Context, name string) (id int, err error) {
	query := "INSERT INTO " + permissionTable + " (name) VALUES(?)"

	err = r.db.QueryRowContext(ctx, query, name).Err()
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	return id, nil
}

func (r *permissionRepository) Update(ctx context.Context, permission domain.Permission) (err error) {
	query := "UPDATE " + permissionTable + " SET name = ? WHERE id = ?"

	err = r.db.QueryRow(query, permission.Name, permission.Id).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *permissionRepository) GetAll(ctx context.Context, page, limit int) (permissions []domain.Permission, err error) {
	query := "SELECT id, name FROM " + permissionTable

	if limit <= 0 {
		limit = defaultLimit
	}

	if page > 0 {
		query = query + " offset " + strconv.Itoa(page) + " limit " + strconv.Itoa(limit)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []domain.Permission{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		permission := domain.Permission{}

		if err = rows.Scan(&permission.Id, &permission.Name); err != nil {
			return []domain.Permission{}, errs.AddTrace(err)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *permissionRepository) GetByName(name string) (permission domain.Permission, err error) {
	query := "SELECT id, name FROM " + permissionTable + " where name = ? limit 1"

	if err = r.db.QueryRow(query, name).Scan(
		&permission.Id,
		&permission.Name,
	); err != nil {
		return domain.Permission{}, errs.AddTrace(err)
	}

	return permission, nil
}

func (r *permissionRepository) GetByRoleId(roleId int) (permissions []domain.Permission, err error) {
	query := "SELECT p.id, p.name FROM " + permissionTable + " as p"
	query = query + " JOIN role_permission rp ON rp.permission_id = p.id"
	query = query + " WHERE rp.role_id = ?"

	rows, err := r.db.Query(query, roleId)
	if err != nil {
		return []domain.Permission{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var permission domain.Permission

		if err = rows.Scan(&permission.Id, &permission.Name); err != nil {
			return []domain.Permission{}, errs.AddTrace(err)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *permissionRepository) GetNamesByUserId(userId int) (permissions []string, err error) {
	query := "SELECT p.name FROM " + permissionTable + " as p"
	query = query + " JOIN role_permission rp ON rp.permission_id = p.id"
	query = query + " JOIN user_role ur ON ur.role_id = rp.role_id"
	query = query + " WHERE ur.user_id = ?"
	permissions, err = r.queryRowsNames(query, userId)
	if err != nil {
		return permissions, errs.AddTrace(err)
	}
	return permissions, nil
}

func (r *permissionRepository) GetNamesByRoleId(roleId int) (permissions []string, err error) {
	query := "SELECT p.name FROM " + permissionTable + " as p"
	query = query + " JOIN role_permission rp ON rp.permission_id = p.id"
	query = query + " WHERE rp.role_id = ?"
	permissions, err = r.queryRowsNames(query, roleId)
	if err != nil {
		return permissions, errs.AddTrace(err)
	}
	return permissions, nil
}

func (r *permissionRepository) queryRowsNames(query string, param int) (permissions []string, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []string{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		if err = rows.Scan(&name); err != nil {
			return []string{}, errs.AddTrace(err)
		}

		permissions = append(permissions, name)
	}

	return permissions, nil
}

func (r *permissionRepository) Delete(ctx context.Context, permissionId int) (err error) {
	query := "DELETE FROM " + permissionTable + " WHERE id = ?"
	err = r.db.QueryRowContext(ctx, query, permissionId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
