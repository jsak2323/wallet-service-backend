package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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

func (r *rpcRequestRepository) Create(rpcRequest rr.RpcRequest) error {
	err := r.db.QueryRow(`
	INSERT INTO `+rpcRequestTable+`(
		arg_name,
		type,
		arg_order,
		source,
		runtime_var_name,
		value,
		rpc_method_id
	)
	VALUES (?,?,?,?,?,?,?);
	`,
		rpcRequest.ArgName,
		rpcRequest.Type,
		rpcRequest.ArgOrder,
		rpcRequest.Source,
		rpcRequest.RuntimeVarName,
		rpcRequest.Value,
		rpcRequest.RpcMethodId,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcRequestRepository) GetByRpcMethodId(rpcMethodId int) (rpcRequests []rr.RpcRequest, err error) {
	query := `SELECT id, arg_name, type, arg_order, source, runtime_var_name, value, rpc_method_id value FROM ` + rpcRequestTable + ` WHERE rpc_method_id = ? ORDER BY arg_order`

	rows, err := r.db.Query(query, rpcMethodId)
	if err != nil {
		return []rr.RpcRequest{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rpcRequest rr.RpcRequest
		err = mapRpcRequest(rows, &rpcRequest)
		if err != nil {
			return []rr.RpcRequest{}, errs.AddTrace(err)
		}

		rpcRequests = append(rpcRequests, rpcRequest)
	}

	return rpcRequests, nil
}

func (r *rpcRequestRepository) Update(rpcRequest rr.RpcRequest) error {
	err := r.db.QueryRow(`
	UPDATE `+rpcRequestTable+`
	SET 
		arg_name = ?,
		type = ?,
		arg_order = ?,
		source = ?,
		runtime_var_name = ?,
		value = ?,
		rpc_method_id = ?
	WHERE id = ?`,
		rpcRequest.ArgName,
		rpcRequest.Type,
		rpcRequest.ArgOrder,
		rpcRequest.Source,
		rpcRequest.RuntimeVarName,
		rpcRequest.Value,
		rpcRequest.RpcMethodId,
		rpcRequest.Id,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
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
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcRequestRepository) Delete(Id int) (err error) {

	query := "DELETE FROM " + rpcRequestTable + " WHERE id = ?"
	err = r.db.QueryRow(query, Id).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
