package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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

func (r *rpcResponseRepository) Create(rpcResponse rr.RpcResponse) error {

	err := r.db.QueryRow(`
	INSERT INTO `+rpcResponseTable+`(
		xml_path,
		field_name,
		data_type_tag,
		rpc_method_id
	)
	VALUES (?,?,?,?);
	`,
		rpcResponse.XMLPath,
		rpcResponse.TargetFieldName,
		rpcResponse.DataTypeXMLTag,
		rpcResponse.ParseType,
		rpcResponse.JsonFieldsStr,
		rpcResponse.RpcMethodId,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcResponseRepository) GetByRpcMethodId(rpcMethodId int) (rpcResponses []rr.RpcResponse, err error) {
	query := `SELECT id, xml_path, field_name, data_type_tag, parse_type, json_fields, rpc_method_id FROM ` + rpcResponseTable + ` WHERE rpc_method_id = ?`

	rows, err := r.db.Query(query, rpcMethodId)
	if err != nil {
		return []rr.RpcResponse{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rpcResponse rr.RpcResponse
		err = mapRpcResponse(rows, &rpcResponse)
		if err != nil {
			return []rr.RpcResponse{}, errs.AddTrace(err)
		}

		rpcResponses = append(rpcResponses, rpcResponse)
	}

	return rpcResponses, nil
}

func (r *rpcResponseRepository) Update(rpcResponse rr.RpcResponse) error {

	err := r.db.QueryRow(`
	UPDATE `+rpcResponseTable+`
	SET 
		xml_path = ?,
		field_name = ?,
		data_type_tag = ?,
		rpc_method_id = ?
	WHERE id = ?`,
		rpcResponse.XMLPath,
		rpcResponse.TargetFieldName,
		rpcResponse.DataTypeXMLTag,
		rpcResponse.ParseType,
		rpcResponse.JsonFieldsStr,
		rpcResponse.RpcMethodId,
		rpcResponse.Id,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func mapRpcResponse(rows *sql.Rows, rpcResponse *rr.RpcResponse) error {
	err := rows.Scan(
		&rpcResponse.Id,
		&rpcResponse.XMLPath,
		&rpcResponse.TargetFieldName,
		&rpcResponse.DataTypeXMLTag,
		&rpcResponse.ParseType,
		&rpcResponse.JsonFieldsStr,
		&rpcResponse.RpcMethodId,
	)
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *rpcResponseRepository) Delete(Id int) (err error) {
	query := "DELETE FROM " + rpcResponseTable + " WHERE id = ?"

	err = r.db.QueryRow(query, Id).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
