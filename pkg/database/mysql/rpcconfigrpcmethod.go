package mysql

import (
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
)

const rpcConfigRpcMethodTable = "rpc_config_rpc_method"

type rpcConfigRpcMethodRepository struct {
	db *sql.DB
}

func NewMysqlRpcConfigRpcMethodRepository(db *sql.DB) domain.Repository {
	return &rpcConfigRpcMethodRepository{
		db,
	}
}

func (r *rpcConfigRpcMethodRepository) Create(rpcConfigId, rpcMethodId int) (err error) {
	query := "INSERT INTO " + rpcConfigRpcMethodTable + " (rpc_config_id, rpc_method_id) VALUES(?, ?)"

	return r.db.QueryRow(query, rpcConfigId, rpcMethodId).Err()
}

func (r *rpcConfigRpcMethodRepository) GetByRpcConfig(rpcConfigId int) (rps []domain.RpcConfigRpcMethod, err error) {
	return r.queryRows("SELECT rpc_config_id, rpc_method_id FROM "+rpcConfigRpcMethodTable+" WHERE rpc_config_id = ?", rpcConfigId)
}

func (r *rpcConfigRpcMethodRepository) GetByRpcMethod(rpcMethodId int) (rps []domain.RpcConfigRpcMethod, err error) {
	return r.queryRows("SELECT rpc_config_id, rpc_method_id FROM "+rpcConfigRpcMethodTable+" WHERE rpc_method_id = ?", rpcMethodId)
}

func (r *rpcConfigRpcMethodRepository) queryRows(query string, param int) (rps []domain.RpcConfigRpcMethod, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.RpcConfigRpcMethod{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.RpcConfigRpcMethod

		if err = rows.Scan(
			&rp.RpcConfigId,
			&rp.RpcMethodId,
		); err != nil {
			return []domain.RpcConfigRpcMethod{}, err
		}

		rps = append(rps, rp)
	}

	return rps, nil
}

func (r *rpcConfigRpcMethodRepository) DeleteByRpcConfig(rpcConfigId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_config_id = ?"

	return r.db.QueryRow(query, rpcConfigId).Err()
}

func (r *rpcConfigRpcMethodRepository) DeleteByRpcMethod(rpcMethodId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_method_id = ?"

	return r.db.QueryRow(query, rpcMethodId).Err()
}

func (r *rpcConfigRpcMethodRepository) Delete(rpcConfigId, rpcMethodId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_config_id = ? and rpc_method_id = ?"

	return r.db.QueryRow(query, rpcConfigId, rpcMethodId).Err()
}