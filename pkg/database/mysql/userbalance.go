package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
	query := "SELECT COALESCE(SUM(" + coin + "), 0) total_" + coin + ", COALESCE(SUM(frozen_" + coin + "), 0) total_frozen_" + coin + " FROM " + userBalanceTable

	if err = r.db.QueryRow(query).Scan(&tcb.Total, &tcb.TotalFrozen); err != nil {
		return domain.TotalCoinBalance{}, errs.AddTrace(err)
	}

	return tcb, nil
}
