package mysql

import (
    "strconv"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
)

const rpcConfigTable = "rpc_config"

type rpcConfigRepository struct {
    db *sql.DB
}

func NewMysqlRpcConfigRepository(db *sql.DB) rc.RpcConfigRepository {
    return &rpcConfigRepository{
        db,
    }
}

func (r *rpcConfigRepository) GetByCurrencyId(currency_id int) ([]rc.RpcConfig, error) {
    query := "SELECT * FROM "+rpcConfigTable
    query += " WHERE currency_id = "+strconv.Itoa(currency_id)
    rpcConfigs := []rc.rpcConfig{}
    var rpcConf rc.RpcConfig

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return rpcConfigs, err }

    for rows.Next() {
        err = rows.Scan(&rpcConf)
        if err != nil { return rpcConfigs, err }

        rpcConfigs = append(rpcConfigs, rpcConf)
    }

    return rpcConfigs, nil
}