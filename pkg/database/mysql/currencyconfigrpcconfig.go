package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const currencyRpcTable = "currency_config_rpc_config"

type currencyRpcRepository struct {
	db *sql.DB
}

func NewMysqlCurrencyRpcRepository(db *sql.DB) domain.Repository {
	return &currencyRpcRepository{
		db,
	}
}

func (r *currencyRpcRepository) Create(ctx context.Context, currencyConfigId, rpcConfigId int) (err error) {
	query := "INSERT INTO " + currencyRpcTable + " (currency_config_id, rpc_config_id) VALUES(?, ?)"

	err = r.db.QueryRowContext(ctx, query, currencyConfigId, rpcConfigId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *currencyRpcRepository) GetByCurrencyConfig(currencyConfigId int) (rps []domain.CurrencyRpc, err error) {
	rps, err = r.queryRows("SELECT currency_config_id, rpc_config_id FROM "+currencyRpcTable+" WHERE currency_config_id = ?", currencyConfigId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *currencyRpcRepository) GetByRpcConfig(rpcConfigId int) (rps []domain.CurrencyRpc, err error) {
	rps, err = r.queryRows("SELECT currency_config_id, rpc_config_id FROM "+currencyRpcTable+" WHERE rpc_config_id = ?", rpcConfigId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *currencyRpcRepository) queryRows(query string, param int) (rps []domain.CurrencyRpc, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.CurrencyRpc{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.CurrencyRpc

		if err = rows.Scan(
			&rp.CurrencyConfigId,
			&rp.RpcConfigId,
		); err != nil {
			return []domain.CurrencyRpc{}, errs.AddTrace(err)
		}

		rps = append(rps, rp)
	}

	return rps, nil
}

func (r *currencyRpcRepository) DeleteByCurrencyConfigId(currencyConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE currency_config_id = ?"
	err = r.db.QueryRow(query, currencyConfigId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *currencyRpcRepository) DeleteByRpcConfigId(rpcConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE rpc_config_id = ?"
	err = r.db.QueryRow(query, rpcConfigId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *currencyRpcRepository) Delete(ctx context.Context, currencyConfigId, rpcConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE currency_config_id = ? and rpc_config_id = ?"
	err = r.db.QueryRowContext(ctx, query, currencyConfigId, rpcConfigId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
