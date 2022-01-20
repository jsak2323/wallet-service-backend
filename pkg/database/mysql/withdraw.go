package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const withdrawTable = "withdraw"
const withdrawDefaultLimit = 10

type withdrawRepository struct {
	db *sql.DB
}

func NewMysqlWithdrawtRepository(db *sql.DB) domain.Repository {
	return &withdrawRepository{
		db,
	}
}

func (r *withdrawRepository) CreateOrUpdate(withdraw domain.Withdraw) (id int, err error) {
	query := `
		INSERT INTO ` + withdrawTable + ` (currency_id, tx, address_to, memo, amount, log_index, confirmations, blockchain_fee, market_price, last_updated)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, now()) ON DUPLICATE KEY UPDATE confirmations = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, errs.AddTrace(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		withdraw.CurrencyId,
		withdraw.Tx,
		withdraw.AddressTo,
		withdraw.Memo,
		withdraw.Amount,
		withdraw.LogIndex,
		withdraw.Confirmations,
		withdraw.BlockchainFee,
		withdraw.MarketPrice,
		withdraw.Confirmations,
	)
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	return int(lastInsertId), nil
}

func (r *withdrawRepository) Get(page, limit int, filters []map[string]interface{}) ([]domain.Withdraw, error) {
	var params []interface{}
	var query string

	query = "SELECT id, currency_id, tx, address_to, memo, amount, log_index, confirmations, blockchain_fee, market_price, success_time, last_updated FROM " + withdrawTable

	if err := parseFilters(filters, &query, &params); err != nil {
		return []domain.Withdraw{}, errs.AddTrace(err)
	}

	if limit <= 0 {
		limit = withdrawDefaultLimit
	}

	if page > 0 {
		page = page * limit
	}

	query = query + " limit ?, ? "
	params = append(params, page, limit)

	return r.queryRows(query, params...)
}

func (r *withdrawRepository) GetById(id int) (withdraw domain.Withdraw, err error) {
	var successTime sql.NullTime

	query := "SELECT id, currency_id, tx, address_to, memo, amount, log_index, confirmations, blockchain_fee, market_price, success_time, last_updated FROM " + withdrawTable + " where id = ?"

	if err = r.db.QueryRow(query, id).Scan(
		&withdraw.Id,
		&withdraw.CurrencyId,
		&withdraw.Tx,
		&withdraw.AddressTo,
		&withdraw.Memo,
		&withdraw.Amount,
		&withdraw.LogIndex,
		&withdraw.Confirmations,
		&withdraw.BlockchainFee,
		&withdraw.MarketPrice,
		&successTime,
		&withdraw.LastUpdated,
	); err != nil {
		return domain.Withdraw{}, errs.AddTrace(err)
	}

	if successTime.Valid {
		withdraw.SuccessTime = successTime.Time.String()
	}

	return withdraw, nil
}

func (r *withdrawRepository) queryRows(query string, params ...interface{}) (deposits []domain.Withdraw, err error) {
	var successTime sql.NullString

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return []domain.Withdraw{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw domain.Withdraw

		if err = rows.Scan(
			&withdraw.Id,
			&withdraw.CurrencyId,
			&withdraw.Tx,
			&withdraw.AddressTo,
			&withdraw.Memo,
			&withdraw.Amount,
			&withdraw.LogIndex,
			&withdraw.Confirmations,
			&withdraw.BlockchainFee,
			&withdraw.MarketPrice,
			&successTime,
			&withdraw.LastUpdated,
		); err != nil {
			return []domain.Withdraw{}, errs.AddTrace(err)
		}

		if successTime.Valid {
			withdraw.SuccessTime = successTime.String
		}

		deposits = append(deposits, withdraw)
	}

	return deposits, nil
}

func (r *withdrawRepository) Update(withdraw domain.Withdraw) (err error) {
	if err = r.db.QueryRow(`
        UPDATE `+withdrawTable+`
        SET 
            currency_id = ?,
            tx = ?,
            address_to = ?,
            memo = ?,
            amount = ?,
            log_index = ?,
            confirmations = ?,
            blockchain_fee = ?,
            market_price = ?,
            success_time = ?,
            last_updated = now()
        WHERE id = ?`,
		withdraw.CurrencyId,
		withdraw.Tx,
		withdraw.AddressTo,
		withdraw.Memo,
		withdraw.Amount,
		withdraw.LogIndex,
		withdraw.Confirmations,
		withdraw.BlockchainFee,
		withdraw.MarketPrice,
		withdraw.SuccessTime,
		withdraw.Id,
	).Err(); err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
