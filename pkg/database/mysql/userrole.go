package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const userRoleTable = "user_role"

type userRoleRepository struct {
	db *sql.DB
}

func NewMysqlUserRoleRepository(db *sql.DB) domain.Repository {
	return &userRoleRepository{
		db,
	}
}

func (r *userRoleRepository) Create(userId, roleId int) (err error) {
	query := "INSERT INTO " + userRoleTable + " (user_id, role_id) VALUES(?,?)"
	if err != nil {
		return errs.AddTrace(err)
	}
	return r.db.QueryRow(query, userId, roleId).Err()
}

func (r *userRoleRepository) GetByUser(userId int) (ur []domain.UserRole, err error) {
	ur, err = r.queryRows("SELECT user_id, role_id FROM "+userRoleTable+" WHERE user_id = ?", userId)
	if err != nil {
		return ur, errs.AddTrace(err)
	}
	return ur, nil
}

func (r *userRoleRepository) GetByRole(roleId int) (ur []domain.UserRole, err error) {

	ur, err = r.queryRows("SELECT user_id, role_id FROM "+userRoleTable+" WHERE role_id = ?", roleId)
	if err != nil {
		return ur, errs.AddTrace(err)
	}
	return ur, nil
}

func (r *userRoleRepository) queryRows(query string, param int) (urs []domain.UserRole, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.UserRole{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ur domain.UserRole

		if err = rows.Scan(
			&ur.UserId,
			&ur.RoleId,
		); err != nil {
			return []domain.UserRole{}, errs.AddTrace(err)
		}

		urs = append(urs, ur)
	}

	return urs, nil
}

func (r *userRoleRepository) DeleteByUserId(userId int) (err error) {
	query := "DELETE FROM " + userRoleTable + " WHERE user_id = ?"

	err = r.db.QueryRow(query, userId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *userRoleRepository) DeleteByRoleId(roleId int) (err error) {
	query := "DELETE FROM " + userRoleTable + " WHERE role_id = ?"

	err = r.db.QueryRow(query, roleId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *userRoleRepository) Delete(userId, roleId int) (err error) {
	query := "DELETE FROM " + userRoleTable + " WHERE user_id = ? and role_id = ?"

	err = r.db.QueryRow(query, userId, roleId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
