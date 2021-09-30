package mysql

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
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

func (r *rpcConfigRepository) GetById(id int) (rc.RpcConfig, error) {
    query := `
        SELECT
            id,
            currency_id,
            type,
            name,
            platform,
            host,
            port,
            path,
            user,
            password,
            hashkey,
            node_version,
            node_last_updated,
            is_health_check_enabled,
            atom_feed,
            address
        FROM ` + rpcConfigTable
    query += " WHERE id = ?"
    rpcConfig := rc.RpcConfig{}

    rows, err := r.db.Query(query, id)
    if err != nil { return rc.RpcConfig{}, err }
    defer rows.Close()

    for rows.Next() {
        err = mapRpcConfig(rows, &rpcConfig)
        if err != nil { return rc.RpcConfig{}, err }
    }

    return rpcConfig, nil
}

func (r *rpcConfigRepository) GetByCurrencyId(currency_id int) ([]rc.RpcConfig, error) {
    query := `
        SELECT
            id,
            currency_id,
            type,
            name,
            platform,
            host,
            port,
            path,
            user,
            password,
            hashkey,
            node_version,
            node_last_updated,
            is_health_check_enabled,
            atom_feed,
            address
        FROM ` + rpcConfigTable
    query += " WHERE currency_id = "+strconv.Itoa(currency_id)
    rpcConfigs := []rc.RpcConfig{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return rpcConfigs, err }

    for rows.Next() {
        var rpcConf rc.RpcConfig
        err = mapRpcConfig(rows, &rpcConf)
        if err != nil { return rpcConfigs, err }

        rpcConfigs = append(rpcConfigs, rpcConf)
    }

    return rpcConfigs, nil
}

func (r *rpcConfigRepository) GetByCurrencySymbol(symbol string) ([]rc.RpcConfig, error) {
    query := `
        SELECT
            `+rpcConfigTable+`.id,
            currency_id,
            type,
            `+rpcConfigTable+`.name,
            platform,
            host,
            port,
            path,
            user,
            password,
            hashkey,
            node_version,
            node_last_updated,
            is_health_check_enabled,
            atom_feed,
            address
        FROM ` + rpcConfigTable
    query += " LEFT JOIN "+currencyConfigTable+" ON "+rpcConfigTable+".currency_id = "+currencyConfigTable+".id "
    query += " WHERE "+currencyConfigTable+".symbol = '"+symbol+"'"
    rpcConfigs := []rc.RpcConfig{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return rpcConfigs, err }

    for rows.Next() {
        var rpcConf rc.RpcConfig
        err = mapRpcConfig(rows, &rpcConf)
        if err != nil { return rpcConfigs, err }

        rpcConfigs = append(rpcConfigs, rpcConf)
    }

    return rpcConfigs, nil
}

func mapRpcConfig(rows *sql.Rows, rpcConf *rc.RpcConfig) error {
    err := rows.Scan(
        &rpcConf.Id,
        &rpcConf.CurrencyId,
        &rpcConf.Type,
        &rpcConf.Name,
        &rpcConf.Platform,
        &rpcConf.Host,
        &rpcConf.Port,
        &rpcConf.Path,
        &rpcConf.User,
        &rpcConf.Password,
        &rpcConf.Hashkey,
        &rpcConf.NodeVersion,
        &rpcConf.NodeLastUpdated,
        &rpcConf.IsHealthCheckEnabled,
        &rpcConf.AtomFeed,
        &rpcConf.Address,
    )

    if err != nil { return err }
    return nil
}


