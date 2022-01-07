package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
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
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(deposit.CurrencyId, deposit.Tx, deposit.AddressTo, deposit.Memo, deposit.Amount, deposit.LogIndex, deposit.Confirmations, deposit.Confirmations)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func (r *depositRepository) Get(page, limit int, filters []map[string]interface{}) ([]domain.Deposit, error) {
	var params []interface{}
	var query string

	query = "SELECT id, currency_id, tx, address_to, memo, amount, log_index, confirmations, success_time, last_updated FROM " + depositTable

	if err := parseFilters(filters, &query, &params); err != nil {
		return []domain.Deposit{}, err
	}

	if limit <= 0 {
		limit = depositDefaultLimit
	}

	if page > 0 {
		page = page * limit
	}

	query = query + " limit ?, ? "
	params = append(params, page, limit)

	return r.queryRows(query, params...)
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
		return domain.Deposit{}, err
	}

	if successTime.Valid {
		deposit.SuccessTime = successTime.Time.String()
	}

	return deposit, nil
}

func (r *depositRepository) queryRows(query string, params ...interface{}) (deposits []domain.Deposit, err error) {
	var successTime sql.NullTime

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return []domain.Deposit{}, err
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
			return []domain.Deposit{}, err
		}

		if successTime.Valid {
			deposit.SuccessTime = successTime.Time.String()
		}

		deposits = append(deposits, deposit)
	}

	return deposits, nil
}

func (r *depositRepository) Update(deposit domain.Deposit) (err error) {
	return r.db.QueryRow(`
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
}