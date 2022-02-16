package mysql

import (
	"context"
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const depositTable = "deposit"
const depositDefaultLimit = 10

type depositRepository struct {
	db *sql.DB
}

func NewMysqlDepositRepository(db *sql.DB) domain.Repository {
	return &depositRepository{
		db,
	}
}

func (r *depositRepository) CreateOrUpdate(deposit domain.Deposit) (id int, err error) {
	query := `
		INSERT INTO ` + depositTable + ` (currency_id, tx, address_to, memo, amount, log_index, confirmations, last_updated)
		VALUES(?, ?, ?, ?, ?, ?, ?, now()) ON DUPLICATE KEY UPDATE confirmations = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, errs.AddTrace(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(deposit.CurrencyId, deposit.Tx, deposit.AddressTo, deposit.Memo, deposit.Amount, deposit.LogIndex, deposit.Confirmations, deposit.Confirmations)
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	return int(lastInsertId), nil
}

func (r *depositRepository) Get(ctx context.Context, page, limit int, filters []map[string]interface{}) ([]domain.Deposit, error) {
	var params []interface{}
	var query string

	query = "SELECT id, currency_id, tx, address_to, memo, amount, log_index, confirmations, success_time, last_updated FROM " + depositTable

	if err := parseFilters(filters, &query, &params); err != nil {
		return []domain.Deposit{}, errs.AddTrace(err)
	}

	if limit <= 0 {
		limit = depositDefaultLimit
	}

	if page > 0 {
		page = page * limit
	}

	query = query + " limit ?, ? "
	params = append(params, page, limit)

	return r.queryRows(ctx, query, params...)
}

func (r *depositRepository) GetById(id int) (deposit domain.Deposit, err error) {
	var successTime sql.NullTime

	query := "SELECT id, currency_id, tx, address_to, memo, amount, log_index, confirmations, success_time, last_updated FROM " + depositTable + " where id = ?"

	if err = r.db.QueryRow(query, id).Scan(
		&deposit.Id,
		&deposit.CurrencyId,
		&deposit.Tx,
		&deposit.AddressTo,
		&deposit.Memo,
		&deposit.Amount,
		&deposit.LogIndex,
		&deposit.Confirmations,
		&successTime,
		&deposit.LastUpdated,
	); err != nil {
		return domain.Deposit{}, errs.AddTrace(err)
	}

	if successTime.Valid {
		deposit.SuccessTime = successTime.Time.String()
	}

	return deposit, nil
}

func (r *depositRepository) queryRows(ctx context.Context, query string, params ...interface{}) (deposits []domain.Deposit, err error) {
	var successTime sql.NullString

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return []domain.Deposit{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var deposit domain.Deposit

		if err = rows.Scan(
			&deposit.Id,
			&deposit.CurrencyId,
			&deposit.Tx,
			&deposit.AddressTo,
			&deposit.Memo,
			&deposit.Amount,
			&deposit.LogIndex,
			&deposit.Confirmations,
			&successTime,
			&deposit.LastUpdated,
		); err != nil {
			return []domain.Deposit{}, errs.AddTrace(err)
		}

		if successTime.Valid {
			deposit.SuccessTime = successTime.String
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func (r *depositRepository) Update(deposit domain.Deposit) (err error) {

	err = r.db.QueryRow(`
	UPDATE `+depositTable+`
	SET 
		currency_id = ?,
		tx = ?,
		address_to = ?,
		memo = ?,
		amount = ?,
		log_index = ?,
		confirmations = ?,
		success_time = ?,
		last_updated = now()
	WHERE id = ?`,
		deposit.CurrencyId,
		deposit.Tx,
		deposit.AddressTo,
		deposit.Memo,
		deposit.Amount,
		deposit.LogIndex,
		deposit.Confirmations,
		deposit.SuccessTime,
		deposit.Id,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
