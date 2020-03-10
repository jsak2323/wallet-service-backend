package mysql

import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    cc "github.com/btcid/wallet-services-backend/pkg/domain/currencyconfig"
)

const currencyConfigTable = "currency_config"

type currencyConfigRepository struct {
    db *sql.DB
}

func NewMysqlCurrencyConfigRepository(db *sql.DB) cc.CurrencyConfigRepository {
    return &currencyConfigRepository{
        db,
    }
}

func (r *currencyConfigRepository) GetAll() ([]cc.CurrencyConfig, error) {
    query := "SELECT * FROM "+currencyConfigTable
    currencyConfigs := []cc.CurrencyConfig{}
    var currConf cc.CurrencyConfig

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return currencyConfigs, err }

    for rows.Next() {
        err = rows.Scan(&currConf)
        if err != nil { return currencyConfigs, err }

        currencyConfigs = append(currencyConfigs, currConf)
    }

    return currencyConfigs, nil
}

func (r *currencyConfigRepository) GetBySymbol(symbol string) (*cc.CurrencyConfig, error) {
    query := "SELECT * FROM "+currencyConfigTable
    query += " WHERE symbol = \""+symbol+"\" "
    var currConf cc.CurrencyConfig

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return &currConf, err }

    for rows.Next() {
        err = rows.Scan(&currConf)
        if err != nil { return &currConf, err }
    }

    return &currConf, nil
}