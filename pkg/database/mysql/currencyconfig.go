package mysql

import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
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

    rows, err := r.db.Query(query)
    defer rows.Close()
    if err != nil { return currencyConfigs, err }

    for rows.Next() { 
        var currConf cc.CurrencyConfig
        err = mapCurrencyConfig(rows, &currConf)
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
        err = mapCurrencyConfig(rows, &currConf)
        if err != nil { return &currConf, err }
    }

    return &currConf, nil
}

func mapCurrencyConfig(rows *sql.Rows, currConf *cc.CurrencyConfig) error {
    var qrCodePrefix    sql.NullString
    var cmcId           sql.NullInt64
    var parentSymbol    sql.NullString

    err := rows.Scan(
        &currConf.Id,
        &currConf.Symbol,
        &currConf.Name,
        &currConf.Unit,
        &currConf.TokenType,
        &currConf.IsFinanceEnabled,
        &currConf.IsSingleAddress,
        &currConf.IsUsingMemo,
        &currConf.IsQrCodeEnabled,
        &currConf.IsAddressNoticeEnabled,
        &qrCodePrefix,
        &currConf.WithdrawFee,
        &currConf.HealthyBlockDiff,
        &currConf.DefaultIdrPrice,
        &cmcId,
        &parentSymbol,
        &currConf.LastUpdated,
    )
    if err != nil { return err }

    if qrCodePrefix.Valid { currConf.QrCodePrefix = qrCodePrefix.String }
    if cmcId.Valid { currConf.CmcId = int(cmcId.Int64) }
    if parentSymbol.Valid { currConf.ParentSymbol = parentSymbol.String }

    return nil
}


