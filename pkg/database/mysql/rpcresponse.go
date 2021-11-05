package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

const rpcResponseTable = "rpc_response"

type rpcResponseRepository struct {
	db *sql.DB
}

func NewMysqlRpcResponseRepository(db *sql.DB) rr.Repository {
	return &rpcResponseRepository{
		db,
	}
}

func (r *rpcResponseRepository) GetByRpcMethodId(rpcConfigId int) (rpcResponses []rr.RpcResponse, err error) {
	query := `SELECT id, xml_path, field_name, data_type_tag, rpc_method_id FROM ` + rpcResponseTable + ` WHERE rpc_method_id = ?`

	rows, err := r.db.Query(query, rpcConfigId)
	if err != nil {
		return []rr.RpcResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpcResponse rr.RpcResponse
		err = mapRpcResponse(rows, &rpcResponse)
		if err != nil {
			return []rr.RpcResponse{}, err
		}

		rpcResponses = append(rpcResponses, rpcResponse)
	}

	return rpcResponses, nil
}

func mapRpcResponse(rows *sql.Rows, rpcResponse *rr.RpcResponse) error {
	err := rows.Scan(
		&rpcResponse.Id,
		&rpcResponse.XMLPath,
		&rpcResponse.FieldName,
		&rpcResponse.DataTypeTag,
		&rpcResponse.RpcMethodId,
	)

	if err != nil {
		return err
	}
	return nil
}
