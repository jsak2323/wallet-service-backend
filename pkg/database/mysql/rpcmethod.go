package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const rpcMethodTable = "rpc_method"

type rpcMethodRepository struct {
	db *sql.DB
}

func NewMysqlRpcMethodRepository(db *sql.DB) rm.Repository {
	return &rpcMethodRepository{
		db,
	}
}

func (r *rpcMethodRepository) Create(rpcMethod rm.RpcMethod) (int, error) {
	stmt, err := r.db.Prepare(`
        INSERT INTO ` + rpcMethodTable + `(
            name,
            type,
            num_of_args,
            network
		)
        VALUES (?,?,?,?);
        `,
	)
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	res, err := stmt.Exec(
		rpcMethod.Name,
		rpcMethod.Type,
		rpcMethod.NumOfArgs,
		rpcMethod.Network,
	)
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, errs.AddTrace(err)
	}

	return int(lastId), nil
}

func (r *rpcMethodRepository) GetAll(page, limit int) (rpcMethods []rm.RpcMethod, err error) {
	query := `
        SELECT
            id,
            name,
            type,
            num_of_args,
            network
        FROM ` + rpcMethodTable

	rows, err := r.db.Query(query)
	if err != nil {
		return []rm.RpcMethod{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		rpcMethod := rm.RpcMethod{}

		if err = mapRpcMethod(rows, &rpcMethod); err != nil {
			return []rm.RpcMethod{}, errs.AddTrace(err)
		}

		rpcMethods = append(rpcMethods, rpcMethod)
	}

	return rpcMethods, nil
}

func (r *rpcMethodRepository) GetByRpcConfigId(rpcConfigId int) (rpcMethods []rm.RpcMethod, err error) {
	query := `
		SELECT 
			rm.id,
			rm.name,
			rm.type,
			rm.num_of_args,
			rm.network
		FROM ` + rpcMethodTable + ` rm
		JOIN ` + rpcConfigRpcMethodTable + ` rcrm ON rcrm.rpc_method_id = rm.id 
		JOIN ` + rpcConfigTable + ` rc ON rcrm.rpc_config_id = rc.id 
		WHERE rcrm.rpc_config_id = ?
	`

	rows, err := r.db.Query(query, rpcConfigId)
	if err != nil {
		return []rm.RpcMethod{}, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rpcMethod = rm.RpcMethod{RpcConfigId: rpcConfigId}
		err = mapRpcMethod(rows, &rpcMethod)
		if err != nil {
			return []rm.RpcMethod{}, errs.AddTrace(err)
		}

		rpcMethods = append(rpcMethods, rpcMethod)
	}

	return rpcMethods, nil
}

func (r *rpcMethodRepository) Update(rpcMethod rm.UpdateRpcMethod) error {
	err := r.db.QueryRow(`
	UPDATE `+rpcMethodTable+`
	SET 
		name = ?,
		type = ?,
		num_of_args = ?,
		network = ?
	WHERE id = ?`,
		rpcMethod.Name,
		rpcMethod.Type,
		rpcMethod.NumOfArgs,
		rpcMethod.Network,
		rpcMethod.Id,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func mapRpcMethod(rows *sql.Rows, rpcMethod *rm.RpcMethod) error {
	err := rows.Scan(
		&rpcMethod.Id,
		&rpcMethod.Name,
		&rpcMethod.Type,
		&rpcMethod.NumOfArgs,
		&rpcMethod.Network,
	)

	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}

func (r *rpcMethodRepository) Delete(Id int) (err error) {
	query := "DELETE FROM " + rpcMethodTable + " WHERE id = ?"
	err = r.db.QueryRow(query, Id).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return nil
}
