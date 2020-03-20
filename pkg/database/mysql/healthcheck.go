package mysql

import (
    "strconv"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

const healthCheckTable = "health_check"

type healthCheckRepository struct {
    db *sql.DB
}

func NewMysqlHealthCheckRepository(db *sql.DB) hc.HealthCheckRepository {
    return &healthCheckRepository{
        db,
    }
}

func (r *healthCheckRepository) GetAll() ([]hc.HealthCheck, error) {
    query := "SELECT * FROM "+healthCheckTable
    allHealthCheck := []hc.HealthCheck{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return allHealthCheck, err }

    for rows.Next() { 
        var healthCheck hc.HealthCheck
        err = mapHealthCheck(rows, &healthCheck)
        if err != nil { return allHealthCheck, err }

        allHealthCheck = append(allHealthCheck, healthCheck)
    }

    return allHealthCheck, nil
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

func (r *healthCheckRepository) Create(healthCheck *hc.HealthCheck) (error) {
    rows, err := r.db.Prepare("INSERT INTO "+healthCheckTable+
        "(rpc_config_id, blockcount, confirm_blockcount) "+
        " VALUES(?, ?, ?)")
    defer rows.Close()
    if err != nil { return err }

    res, err := rows.Exec(
        healthCheck.RpcConfigId, healthCheck.BlockCount, healthCheck.ConfirmBlockCount)
    if err != nil { return err }

    lastInsertId, _ := res.LastInsertId()
    healthCheck.Id = int(lastInsertId)

    return nil
}

func (r *healthCheckRepository) Update(id int, healthCheck hc.HealthCheck) (error) {
    query := "UPDATE "+healthCheckTable+" SET "+
    " `rpc_config_id` = "+strconv.Itoa(healthCheck.RpcConfigId)+", "+
    " `blockcount` = "+strconv.Itoa(healthCheck.BlockCount)+", "+
    " `confirm_blockcount` = "+strconv.Itoa(healthCheck.ConfirmBlockCount)+", "+
    " `last_updated_block_num` = CURRENT_TIMESTAMP() "+
    " WHERE id = "+strconv.Itoa(id)

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
        &healthCheck.ConfirmBlockCount,
        &healthCheck.LastUpdated,
    )
    if err != nil { return err }
    return nil
}