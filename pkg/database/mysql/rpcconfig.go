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

func NewMysqlRpcConfigRepository(db *sql.DB) rc.Repository {
	return &rpcConfigRepository{
		db,
	}
}

func (r *rpcConfigRepository) Create(rpcConfig rc.RpcConfig) error {
	return r.db.QueryRow(`
        INSERT INTO `+rpcConfigTable+`(
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
			address)
        VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
        `,
		rpcConfig.Type,
		rpcConfig.Name,
		rpcConfig.Platform,
		rpcConfig.Host,
		rpcConfig.Port,
		rpcConfig.Path,
		rpcConfig.User,
		rpcConfig.Password,
		rpcConfig.Hashkey,
		rpcConfig.NodeVersion,
		rpcConfig.NodeLastUpdated,
		rpcConfig.IsHealthCheckEnabled,
		rpcConfig.AtomFeed,
		rpcConfig.Address,
	).Err()
}

func (r *rpcConfigRepository) GetAll(page, limit int) (rpcConfigs []rc.RpcConfig, err error) {
	query := `
        SELECT
            id,
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
            address,
			active
        FROM ` + rpcConfigTable

	rows, err := r.db.Query(query)
	if err != nil {
		return []rc.RpcConfig{}, err
	}
	defer rows.Close()

	for rows.Next() {
		rpcConfig := rc.RpcConfig{}

		if err = mapRpcConfig(rows, &rpcConfig); err != nil {
			return []rc.RpcConfig{}, err
		}

		rpcConfigs = append(rpcConfigs, rpcConfig)
	}

	return rpcConfigs, nil
}

func (r *rpcConfigRepository) GetById(id int) (rc.RpcConfig, error) {
	query := `
        SELECT
            id,
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
            address,
			active
        FROM ` + rpcConfigTable + ` WHERE id = ?`

	rpcConfig := rc.RpcConfig{}

	rows, err := r.db.Query(query, id)
	if err != nil {
		return rc.RpcConfig{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = mapRpcConfig(rows, &rpcConfig)
		if err != nil {
			return rc.RpcConfig{}, err
		}
	}

	return rpcConfig, nil
}

func (r *rpcConfigRepository) GetByCurrencyId(currency_id int) ([]rc.RpcConfig, error) {
	query := `
        SELECT
            id,
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
            address,
			active
        FROM ` + rpcConfigTable + `
		JOIN ` + currencyRpcTable + ` on ` + currencyRpcTable + `.rpc_config_id = ` + rpcConfigTable + `.id
		`
	query += " WHERE "+ currencyRpcTable +".currency_config_id = " + strconv.Itoa(currency_id)
	rpcConfigs := []rc.RpcConfig{}

	rows, err := r.db.Query(query)
	if err != nil {
		return rpcConfigs, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpcConf rc.RpcConfig
		err = mapRpcConfig(rows, &rpcConf)
		if err != nil {
			return rpcConfigs, err
		}

		rpcConfigs = append(rpcConfigs, rpcConf)
	}

	return rpcConfigs, nil
}

func (r *rpcConfigRepository) GetByCurrencySymbol(symbol string) ([]rc.RpcConfig, error) {
	query := `
        SELECT
            ` + rpcConfigTable + `.id,
            type,
            ` + rpcConfigTable + `.name,
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
            address,
            ` + rpcConfigTable + `.active
        FROM ` + rpcConfigTable
	query += " LEFT JOIN " + currencyConfigTable + " ON " + rpcConfigTable + ".currency_id = " + currencyConfigTable + ".id "
	query += " WHERE " + currencyConfigTable + ".symbol = '" + symbol + "'"
	rpcConfigs := []rc.RpcConfig{}

	rows, err := r.db.Query(query)
	if err != nil {
		return rpcConfigs, err
	}
	defer rows.Close()

	for rows.Next() {
		var rpcConf rc.RpcConfig
		err = mapRpcConfig(rows, &rpcConf)
		if err != nil {
			return rpcConfigs, err
		}

		rpcConfigs = append(rpcConfigs, rpcConf)
	}

	return rpcConfigs, nil
}

func mapRpcConfig(rows *sql.Rows, rpcConf *rc.RpcConfig) error {
	err := rows.Scan(
		&rpcConf.Id,
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
		&rpcConf.Active,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *rpcConfigRepository) Update(rpcConfig rc.RpcConfig) (err error) {
	return r.db.QueryRow(`
        UPDATE `+rpcConfigTable+`
        SET 
            type = ?,
            name = ?,
            platform = ?,
            host = ?,
            port = ?,
            path = ?,
            user = ?,
            password = ?,
            hashkey = ?,
            node_version = ?,
            node_last_updated = ?,
            is_health_check_enabled = ?,
            atom_feed = ?,
            address = ?
        WHERE id = ?`,
		rpcConfig.Type,
		rpcConfig.Name,
		rpcConfig.Platform,
		rpcConfig.Host,
		rpcConfig.Port,
		rpcConfig.Path,
		rpcConfig.User,
		rpcConfig.Password,
		rpcConfig.Hashkey,
		rpcConfig.NodeVersion,
		rpcConfig.NodeLastUpdated,
		rpcConfig.IsHealthCheckEnabled,
		rpcConfig.AtomFeed,
		rpcConfig.Address,
		rpcConfig.Id,
	).Err()
}

func (r *rpcConfigRepository) ToggleActive(userId int, active bool) error {
	query := "UPDATE " + rpcConfigTable + " SET active = ? WHERE id = ?"

	return r.db.QueryRow(query, active, userId).Err()
}
