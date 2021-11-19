package mysql

import (
    "strconv"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
)

const healthCheckTable = "health_check"

type healthCheckRepository struct {
    db *sql.DB
}

func NewMysqlHealthCheckRepository(db *sql.DB) hc.Repository {
    return &healthCheckRepository{
        db,
    }
}

func (r *healthCheckRepository) GetAll() ([]hc.HealthCheck, error) {
    query := "SELECT * FROM "+healthCheckTable
    healthChecks := []hc.HealthCheck{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return healthChecks, err }

    for rows.Next() { 
        var healthCheck hc.HealthCheck
        err = mapHealthCheck(rows, &healthCheck)
        if err != nil { return healthChecks, err }

        healthChecks = append(healthChecks, healthCheck)
    }

    return healthChecks, nil
}

func (r *healthCheckRepository) GetAllWithRpcConfig() ([]hc.HealthCheckWithRpcConfig, error) {
    query := "SELECT * FROM "+healthCheckTable+" hct "
    query += "JOIN "+rpcConfigTable+" rct ON rct.id = hct.rpc_config_id "
    healthChecks := []hc.HealthCheckWithRpcConfig{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return healthChecks, err }

    for rows.Next() { 
        var healthCheck hc.HealthCheckWithRpcConfig
        err = mapHealthCheckWithRpcConfig(rows, &healthCheck)
        if err != nil { return healthChecks, err }

        healthChecks = append(healthChecks, healthCheck)
    }

    return healthChecks, nil
}

func (r *healthCheckRepository) GetByRpcConfigId(rpcConfigId int) (hc.HealthCheck, error) {
    query := "SELECT * FROM "+healthCheckTable+" WHERE rpc_config_id = "+strconv.Itoa(rpcConfigId)
    var healthCheck hc.HealthCheck

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return healthCheck, err }

    for rows.Next() { 
        err = mapHealthCheck(rows, &healthCheck)
        if err != nil { return healthCheck, err }
    }

    return healthCheck, nil
}

func (r *healthCheckRepository) GetLastUpdatedTime() (string, error) {
    query := "SELECT * FROM "+healthCheckTable+" ORDER BY last_updated DESC LIMIT 1"
    var healthCheck hc.HealthCheck

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return "", err }

    for rows.Next() { 
        err = mapHealthCheck(rows, &healthCheck)
        if err != nil { return "", err }
    }

    return healthCheck.LastUpdated, nil
}

func (r *healthCheckRepository) Create(healthCheck *hc.HealthCheck) (error) {
    rows, err := r.db.Prepare("INSERT INTO "+healthCheckTable+
        "(rpc_config_id, blockcount, block_diff, is_healthy) "+
        " VALUES(?, ?, ?, ?)")
    defer rows.Close()
    if err != nil { return err }

    res, err := rows.Exec(
        healthCheck.RpcConfigId, 
        healthCheck.BlockCount, 
        healthCheck.BlockDiff, 
        healthCheck.IsHealthy,
    )
    if err != nil { return err }

    lastInsertId, _ := res.LastInsertId()
    healthCheck.Id = int(lastInsertId)

    return nil
}

func (r *healthCheckRepository) Update(healthCheck *hc.HealthCheck) (error) {
    isHealthy := 0
    if healthCheck.IsHealthy {
        isHealthy = 1
    }

    query := "UPDATE "+healthCheckTable+" SET "+
    " `blockcount` = "+strconv.Itoa(healthCheck.BlockCount)+", "+
    " `block_diff` = "+strconv.Itoa(healthCheck.BlockDiff)+", "+
    " `is_healthy` = "+strconv.Itoa(isHealthy)+", "+
    " `last_updated` = CURRENT_TIMESTAMP() "+
    " WHERE id = "+strconv.Itoa(healthCheck.Id)

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return err }

    return nil
}

func mapHealthCheck(rows *sql.Rows, healthCheck *hc.HealthCheck) error {
    err := rows.Scan(
        &healthCheck.Id,
        &healthCheck.RpcConfigId,
        &healthCheck.BlockCount,
        &healthCheck.BlockDiff,
        &healthCheck.IsHealthy,
        &healthCheck.LastUpdated,
    )
    if err != nil { return err }
    return nil
}

func mapHealthCheckWithRpcConfig(rows *sql.Rows, healthCheck *hc.HealthCheckWithRpcConfig) error {
    placeHolderInt := 0
    
    err := rows.Scan(
        &healthCheck.Id,
        &healthCheck.RpcConfig.Id,
        &healthCheck.BlockCount,
        &healthCheck.BlockDiff,
        &healthCheck.IsHealthy,
        &healthCheck.LastUpdated,
        &placeHolderInt,
        &healthCheck.RpcConfig.Type,
        &healthCheck.RpcConfig.Name,
        &healthCheck.RpcConfig.Platform,
        &healthCheck.RpcConfig.Host,
        &healthCheck.RpcConfig.Port,
        &healthCheck.RpcConfig.Path,
        &healthCheck.RpcConfig.User,
        &healthCheck.RpcConfig.Password,
        &healthCheck.RpcConfig.Hashkey,
        &healthCheck.RpcConfig.NodeVersion,
        &healthCheck.RpcConfig.NodeLastUpdated,
        &healthCheck.RpcConfig.IsHealthCheckEnabled,
        &healthCheck.RpcConfig.AtomFeed,
    )

    if err != nil { return err }
    return nil
}


