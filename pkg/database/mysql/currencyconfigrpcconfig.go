package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
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

func (r *currencyRpcRepository) Create(currencyConfigId, rpcConfigId int) (err error) {
	query := "INSERT INTO " + currencyRpcTable + " (currency_config_id, rpc_config_id) VALUES(?, ?)"

	return r.db.QueryRow(query, currencyConfigId, rpcConfigId).Err()
}

func (r *currencyRpcRepository) GetByCurrencyConfig(currencyConfigId int) (rps []domain.CurrencyRpc, err error) {
	return r.queryRows("SELECT currency_config_id, rpc_config_id FROM "+currencyRpcTable+" WHERE currency_config_id = ?", currencyConfigId)
}

func (r *currencyRpcRepository) GetByRpcConfig(rpcConfigId int) (rps []domain.CurrencyRpc, err error) {
	return r.queryRows("SELECT currency_config_id, rpc_config_id FROM "+currencyRpcTable+" WHERE rpc_config_id = ?", rpcConfigId)
}

func (r *currencyRpcRepository) queryRows(query string, param int) (rps []domain.CurrencyRpc, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.CurrencyRpc{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.CurrencyRpc

		if err = rows.Scan(
			&rp.CurrencyConfigId,
			&rp.RpcConfigId,
		); err != nil {
			return []domain.CurrencyRpc{}, err
		}

		rps = append(rps, rp)
	}

	return rps, nil
}

func (r *currencyRpcRepository) DeleteByCurrencyConfigId(currencyConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE currency_config_id = ?"

	return r.db.QueryRow(query, currencyConfigId).Err()
}

func (r *currencyRpcRepository) DeleteByRpcConfigId(rpcConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE rpc_config_id = ?"

	return r.db.QueryRow(query, rpcConfigId).Err()
}

func (r *currencyRpcRepository) Delete(currencyConfigId, rpcConfigId int) (err error) {
	query := "DELETE FROM " + currencyRpcTable + " WHERE currency_config_id = ? and rpc_config_id = ?"

	return r.db.QueryRow(query, currencyConfigId, rpcConfigId).Err()
}
