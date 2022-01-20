package mysql

import (
	"database/sql"
	"strings"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/withdrawexchange"
)

type withdrawMysqlRepository struct {
	db *sql.DB
}

func NewMysqlWithdrawRepository(db *sql.DB) domain.Repository {
	return &withdrawMysqlRepository{
		db,
	}
}

func (r *withdrawMysqlRepository) GetPendingWithdraw(symbol string) (result string, err error) {
	symbol = strings.ToLower(symbol)

	query := "SELECT SUM(amount) as sum_amount FROM withdraw_"+symbol +" WHERE status in ('pending', 'approved')"

	sumAmount := sql.NullString{}

	if err = r.db.QueryRow(query).Scan(&sumAmount); err != nil {
		return "0", err
	}

	if sumAmount.Valid { result = sumAmount.String }
	
	return result, nil
}