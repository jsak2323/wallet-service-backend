package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

type mysqlMarketRepository struct {
	db *sql.DB
}

func NewMysqlMarketRepository(db *sql.DB) domain.Repository {
	return &mysqlMarketRepository{
		db,
	}
}

func (r *mysqlMarketRepository) LastPriceBySymbol(symbol, trade string) (price string, err error) {
	query := "SELECT price FROM "+symbol+trade+"_trades ORDER BY id DESC LIMIT 1"

	if err = r.db.QueryRow(query).Scan(&price); err != nil {
		return "0", err
	}
	
	return price, nil
}
