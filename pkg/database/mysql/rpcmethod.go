package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
)

const rpcMethodTable = "rpc_method"
const rpcConfigRpcMethodTable = "rpc_config_rpc_method"

type rpcMethodRepository struct {
	db *sql.DB
}

func NewMysqlRpcMethodRepository(db *sql.DB) rm.Repository {
	return &rpcMethodRepository{
		db,
	}
}

func (r *rpcMethodRepository) GetByRpcConfigId(rpcConfigId int) (rpcMethods []rm.RpcMethod, err error) {
	query := `
		SELECT 
			rm.id,
			rm.name,
			rm.type,
			rm.num_of_args
		FROM ` + rpcMethodTable + ` rm
		JOIN ` + rpcConfigRpcMethodTable + ` rcrm ON rcrm.rpc_method_id = rm.id 
		JOIN ` + rpcConfigTable + ` rc ON rcrm.rpc_config_id = rc.id 
		WHERE rcrm.rpc_config_id = ?
	`

	rows, err := r.db.Query(query, rpcConfigId)
	if err != nil {
		return []rm.RpcMethod{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpcMethod rm.RpcMethod
		err = mapRpcMethod(rows, &rpcMethod)
		if err != nil {
			return []rm.RpcMethod{}, err
		}

		rpcMethods = append(rpcMethods, rpcMethod)
	}

	return rpcMethods, nil
}

func mapRpcMethod(rows *sql.Rows, rpcMethod *rm.RpcMethod) error {
	err := rows.Scan(
		&rpcMethod.Id,
		&rpcMethod.Name,
		&rpcMethod.Type,
		&rpcMethod.NumOfArgs,
	)

	if err != nil {
		return err
	}
	return nil
}
