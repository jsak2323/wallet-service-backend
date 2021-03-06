package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const systemConfigTable = "system_config"

type systemConfigRepository struct {
	db *sql.DB
}

func NewMysqlSystemConfigRepository(db *sql.DB) sc.Repository {
	return &systemConfigRepository{
		db,
	}
}

func (r *systemConfigRepository) GetAll() ([]sc.SystemConfig, error) {
	query := "SELECT * FROM " + systemConfigTable
	systemConfigs := []sc.SystemConfig{}

	rows, err := r.db.Query(query)
	defer rows.Close()
	if err != nil {
		return systemConfigs, errs.AddTrace(err)
	}

	for rows.Next() {
		var sysConf sc.SystemConfig
		err = mapSystemConfig(rows, &sysConf)
		if err != nil {
			return systemConfigs, errs.AddTrace(err)
		}

		systemConfigs = append(systemConfigs, sysConf)
	}

	return systemConfigs, nil
}

func (r *systemConfigRepository) GetByName(ctx context.Context, configName string) (*sc.SystemConfig, error) {
	query := "SELECT * FROM " + systemConfigTable
	query += " WHERE name = \"" + configName + "\" "
	var sysConf sc.SystemConfig

	rows, err := r.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return &sysConf, errs.AddTrace(err)
	}

	for rows.Next() {
		err = mapSystemConfig(rows, &sysConf)
		if err != nil {
			return &sysConf, errs.AddTrace(err)
		}
	}

	return &sysConf, nil
}

func (r *systemConfigRepository) Update(sysConf sc.SystemConfig) error {
	query := "UPDATE " + systemConfigTable + " SET " +
		" `value` = \"" + sysConf.Value + "\" " +
		" WHERE `name` = \"" + sysConf.Name + "\""

	rows, err := r.db.Query(query)
	defer rows.Close()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func mapSystemConfig(rows *sql.Rows, sysConf *sc.SystemConfig) error {
	err := rows.Scan(
		&sysConf.Name,
		&sysConf.Value,
	)
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
