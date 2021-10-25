package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
)

const rpcRequestTable = "rpc_request"

type rpcRequestRepository struct {
	db *sql.DB
}

func NewMysqlRpcRequestRepository(db *sql.DB) rr.Repository {
	return &rpcRequestRepository{
		db,
	}
}

func (r *rpcRequestRepository) GetByRpcMethodId(rpcConfigId int) (rpcRequests []rr.RpcRequest, err error) {
	query := `SELECT id, arg_name, type, arg_order, source, runtime_var_name, value, rpc_method_id value FROM ` + rpcRequestTable + ` WHERE rpc_method_id = ? ORDER BY arg_order`

	rows, err := r.db.Query(query, rpcConfigId)
	if err != nil {
		return []rr.RpcRequest{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpcRequest rr.RpcRequest
		err = mapRpcRequest(rows, &rpcRequest)
		if err != nil {
			return []rr.RpcRequest{}, err
		}

		rpcRequests = append(rpcRequests, rpcRequest)
	}

	return rpcRequests, nil
}

func mapRpcRequest(rows *sql.Rows, rpcRequest *rr.RpcRequest) error {
	err := rows.Scan(
		&rpcRequest.Id,
		&rpcRequest.ArgName,
		&rpcRequest.Type,
		&rpcRequest.ArgOrder,
		&rpcRequest.Source,
		&rpcRequest.RuntimeVarName,
		&rpcRequest.Value,
		&rpcRequest.RpcMethodId,
	)

	if err != nil {
		return err
	}
	return nil
}
