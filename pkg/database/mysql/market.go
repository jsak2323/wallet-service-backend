package mysql

import (
	"database/sql"
	"errors"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/market"
)

type marketMysqlRepository struct {
	db *sql.DB
}

func NewMarketMysqlRepository(db *sql.DB) domain.Repository {
	return &marketMysqlRepository{
		db,
	}
}

func (r *marketMysqlRepository) LastPriceBySymbol(symbol, trade string) (price string, err error) {
	query := "SELECT price FROM "+symbol+trade+"_trades ORDER BY id DESC LIMIT 1"

	if err = r.db.QueryRow(query).Scan(&price); err != nil {
		return "0", err
	}
	
	return price, nil
}
