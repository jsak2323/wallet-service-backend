package mysql

import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
)

const systemConfigTable = "system_config"

type systemConfigRepository struct {
    db *sql.DB
}

func NewMysqlSystemConfigRepository(db *sql.DB) sc.SystemConfigRepository {
    return &systemConfigRepository{
        db,
    }
}

func (r *systemConfigRepository) GetAll() ([]sc.SystemConfig, error) {
    query := "SELECT * FROM "+systemConfigTable
    systemConfigs := []sc.SystemConfig{}

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return systemConfigs, err }

    for rows.Next() { 
        var sysConf sc.SystemConfig
        err = mapSystemConfig(rows, &sysConf)
        if err != nil { return systemConfigs, err }

        systemConfigs = append(systemConfigs, sysConf)
    }

    return systemConfigs, nil
}

func (r *systemConfigRepository) GetByName(configName string) (*sc.SystemConfig, error) {
    query := "SELECT * FROM "+systemConfigTable
    query += " WHERE name = \""+configName+"\" "
    var sysConf sc.SystemConfig

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return &sysConf, err }

    for rows.Next() {
        err = mapSystemConfig(rows, &sysConf)
        if err != nil { return &sysConf, err }
    }

    return &sysConf, nil
}

func mapSystemConfig(rows *sql.Rows, sysConf *sc.SystemConfig) error {
    err := rows.Scan(
        &sysConf.Name,
        &sysConf.Value,
    )
    if err != nil { return err }

    return nil
}


