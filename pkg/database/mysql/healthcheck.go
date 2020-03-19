package mysql

import (
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