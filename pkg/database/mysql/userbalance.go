package mysql

import (
	"database/sql"
	
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
)

const userBalanceTable = "user_balance"
const userBalanceDefaultLimit = 10

type userBalanceRepository struct {
	db *sql.DB
}

func NewMysqlUserBalanceRepository(db *sql.DB) domain.Repository {
    return &userBalanceRepository{
        db,
    }
}

func (r *userBalanceRepository) GetTotalCoinBalance(coin string) (tcb domain.TotalCoinBalance, err error) {
	query := "SELECT SUM("+coin+") total_"+coin+", SUM(frozen_"+coin+") total_frozen_"+coin+" FROM "+userBalanceTable

	var total sql.NullInt64
	var totalFrozen sql.NullInt64
	
	if err = r.db.QueryRow(query).Scan(&total, &totalFrozen); err != nil {
		return domain.TotalCoinBalance{}, err
	}

	if total.Valid { tcb.Total = total.Int64 }
	if total.Valid { tcb.TotalFrozen = totalFrozen.Int64 }

	return tcb, nil
}