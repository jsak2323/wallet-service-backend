package mysql

import (
	"context"
	"database/sql"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfigrpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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

func (r *rpcConfigRpcMethodRepository) Create(ctx context.Context, rpcConfigId, rpcMethodId int) (err error) {
	query := "INSERT INTO " + rpcConfigRpcMethodTable + " (rpc_config_id, rpc_method_id) VALUES(?, ?)"
	err = r.db.QueryRowContext(ctx, query, rpcConfigId, rpcMethodId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcConfigRpcMethodRepository) GetByRpcConfig(rpcConfigId int) (rps []domain.RpcConfigRpcMethod, err error) {
	rps, err = r.queryRows("SELECT rpc_config_id, rpc_method_id FROM "+rpcConfigRpcMethodTable+" WHERE rpc_config_id = ?", rpcConfigId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *rpcConfigRpcMethodRepository) GetByRpcMethod(rpcMethodId int) (rps []domain.RpcConfigRpcMethod, err error) {
	rps, err = r.queryRows("SELECT rpc_config_id, rpc_method_id FROM "+rpcConfigRpcMethodTable+" WHERE rpc_method_id = ?", rpcMethodId)
	if err != nil {
		return rps, errs.AddTrace(err)
	}
	return rps, nil
}

func (r *rpcConfigRpcMethodRepository) queryRows(query string, param int) (rps []domain.RpcConfigRpcMethod, err error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return []domain.RpcConfigRpcMethod{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rp domain.RpcConfigRpcMethod

		if err = rows.Scan(
			&rp.RpcConfigId,
			&rp.RpcMethodId,
		); err != nil {
			return []domain.RpcConfigRpcMethod{}, errs.AddTrace(err)
		}

		rps = append(rps, rp)
	}

	return rps, nil
}

func (r *rpcConfigRpcMethodRepository) DeleteByRpcConfig(rpcConfigId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_config_id = ?"
	err = r.db.QueryRow(query, rpcConfigId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcConfigRpcMethodRepository) DeleteByRpcMethod(ctx context.Context, rpcMethodId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_method_id = ?"
	err = r.db.QueryRowContext(ctx, query, rpcMethodId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcConfigRpcMethodRepository) Delete(ctx context.Context, rpcConfigId, rpcMethodId int) (err error) {
	query := "DELETE FROM " + rpcConfigRpcMethodTable + " WHERE rpc_config_id = ? and rpc_method_id = ?"
	err = r.db.QueryRowContext(ctx, query, rpcConfigId, rpcMethodId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
