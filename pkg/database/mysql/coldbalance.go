package mysql

import (
	"database/sql"
	
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
)

const coldBalanceTable = "cold_balance"
const coldBalanceDefaultLimit = 10

type coldBalanceRepository struct {
	db *sql.DB
}

func NewMysqlColdBalanceRepository(db *sql.DB) domain.Repository {
    return &coldBalanceRepository{
        db,
    }
}

func (r *coldBalanceRepository) Create(coldBalance domain.ColdBalance) (id int, err error) {
	query := `
		INSERT INTO ` + coldBalanceTable + ` (currency_id, name, type, fireblocks_name, balance, address, last_updated) 
		VALUES(?, ?, ?, ?, ?, ?, now())
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
    defer stmt.Close()

	res, err := stmt.Exec(coldBalance.CurrencyId, coldBalance.Name, coldBalance.Type, coldBalance.FireblocksName, coldBalance.Balance, coldBalance.Address); 
	if err != nil {
        return 0, err
	}

	lastInsertId, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

	return int(lastInsertId), nil
}

func(r *coldBalanceRepository) GetAll(page, limit int) ([]domain.ColdBalance, error) {
	var params []interface{}
	
	query := "SELECT id, currency_id, name, type, fireblocks_name, balance, address, active, last_updated FROM " + coldBalanceTable

	if limit <= 0 {
		limit = coldBalanceDefaultLimit
	}
	
	if page > 0 {
		query = query + " offset ?"
		params = append(params, page*limit)
	}
	
	query = query + " limit ?"
	params = append(params, limit)
	
	return r.queryRows(query, params...)
}

func (r *coldBalanceRepository) GetByName(name string) (balance domain.ColdBalance, err error) {
	query := "SELECT id, currency_id, name, type, fireblocks_name, balance, address, active, last_updated FROM " + coldBalanceTable + " where name = ?"

	if err = r.db.QueryRow(query, name).Scan(
		&balance.Id,
		&balance.CurrencyId,
		&balance.Name,
		&balance.Type,
		&balance.FireblocksName,
		&balance.Balance,
		&balance.Address,
		&balance.Active,
		&balance.LastUpdated,
	); err != nil {
		return domain.ColdBalance{}, err
	}

	return balance, nil
}

func (r *coldBalanceRepository) GetByCurrencyId(currencyId int) (balances []domain.ColdBalance, err error) {
	query := "SELECT id, currency_id, name, type, fireblocks_name, balance, address, active, last_updated FROM " + coldBalanceTable + " where currency_id = ?"
	
	return r.queryRows(query, currencyId)
}

func (r *coldBalanceRepository) queryRows(query string, params... interface{}) (balances []domain.ColdBalance, err error) {
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return []domain.ColdBalance{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var balance domain.ColdBalance
		
		if err = rows.Scan(
			&balance.Id,
			&balance.CurrencyId,
			&balance.Name,
			&balance.Type,
			&balance.FireblocksName,
			&balance.Balance,
			&balance.Address,
			&balance.Active,
			&balance.LastUpdated,
		); err != nil {
			return []domain.ColdBalance{}, err
		}

		balances = append(balances, balance)
	}

	return balances, nil
}

// not updating balance
func (r *coldBalanceRepository) Update(coldBalance domain.ColdBalance) (err error) {
	return r.db.QueryRow(`
        UPDATE ` + coldBalanceTable + `
        SET 
            currency_id = ?,
            name = ?,
            type = ?,
            fireblocks_name = ?,
            address = ?,
            last_updated = now()
        WHERE id = ?`,
        coldBalance.CurrencyId,
        coldBalance.Name,
        coldBalance.Type,
        coldBalance.FireblocksName,
        coldBalance.Address,
        coldBalance.Id,
    ).Err()
}

func (r *coldBalanceRepository) UpdateBalance(id int, balance string) (err error) {
	return r.db.QueryRow("UPDATE cold_balance SET balance = ? WHERE id = ?", balance, id).Err()
}

func(r *coldBalanceRepository) ToggleActive(userId int, active bool) error {
	query := "UPDATE " + coldBalanceTable + " SET active = ? WHERE id = ?"

	return r.db.QueryRow(query, active, userId).Err()
}