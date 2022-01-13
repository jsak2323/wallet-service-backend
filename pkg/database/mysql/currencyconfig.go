package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

const currencyConfigTable = "currency_config"

type currencyConfigRepository struct {
	db *sql.DB
}

func NewMysqlCurrencyConfigRepository(db *sql.DB) cc.Repository {
	return &currencyConfigRepository{
		db,
	}
}

// errs.AddTrace(err)

func (r *currencyConfigRepository) Create(currencyConfig cc.CurrencyConfig) error {

	err := r.db.QueryRow(`
    INSERT INTO `+currencyConfigTable+`(
        symbol,
        name,
        unit,
        token_type,
        is_finance_enabled,
        is_single_address,
        is_using_memo,
        is_qrcode_enabled,
        is_address_notice_enabled,
        qrcode_prefix,
        withdraw_fee,
        healthy_block_diff,
        default_idr_price,
        cmc_id,
        parent_symbol,
        module_type,
        last_updated)
    VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now());
    `,
		currencyConfig.Symbol,
		currencyConfig.Name,
		currencyConfig.Unit,
		currencyConfig.TokenType,
		currencyConfig.IsFinanceEnabled,
		currencyConfig.IsSingleAddress,
		currencyConfig.IsUsingMemo,
		currencyConfig.IsQrCodeEnabled,
		currencyConfig.IsAddressNoticeEnabled,
		currencyConfig.QrCodePrefix,
		currencyConfig.WithdrawFee,
		currencyConfig.HealthyBlockDiff,
		currencyConfig.DefaultIdrPrice,
		currencyConfig.CmcId,
		currencyConfig.ParentSymbol,
		currencyConfig.ModuleType,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func (r *currencyConfigRepository) GetAll() ([]cc.CurrencyConfig, error) {
	query := `
        SELECT 
            id,
            symbol,
            name,
            unit,
            token_type,
            is_finance_enabled,
            is_single_address,
            is_using_memo,
            is_qrcode_enabled,
            is_address_notice_enabled,
            qrcode_prefix,
            withdraw_fee,
            healthy_block_diff,
            default_idr_price,
            cmc_id,
            parent_symbol,
            module_type,
            active,
            last_updated
        FROM ` + currencyConfigTable
	currencyConfigs := []cc.CurrencyConfig{}

	rows, err := r.db.Query(query)
	if err != nil {
		return currencyConfigs, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var currConf cc.CurrencyConfig
		err = mapCurrencyConfig(rows, &currConf)
		if err != nil {
			return currencyConfigs, errs.AddTrace(err)
		}

		currencyConfigs = append(currencyConfigs, currConf)
	}

	return currencyConfigs, nil
}

func (r *currencyConfigRepository) GetBySymbol(symbol string) (*cc.CurrencyConfig, error) {
	query := `
        SELECT 
            id,
            symbol,
            name,
            unit,
            token_type,
            is_finance_enabled,
            is_single_address,
            is_using_memo,
            is_qrcode_enabled,
            is_address_notice_enabled,
            qrcode_prefix,
            withdraw_fee,
            healthy_block_diff,
            default_idr_price,
            cmc_id,
            parent_symbol,
            module_type,
            active,
            last_updated
        FROM ` + currencyConfigTable + `
        WHERE symbol = ?`
	var currConf cc.CurrencyConfig

	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return &currConf, errs.AddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = mapCurrencyConfig(rows, &currConf)
		if err != nil {
			return &currConf, errs.AddTrace(err)
		}
	}

	return &currConf, nil
}

func mapCurrencyConfig(rows *sql.Rows, currConf *cc.CurrencyConfig) error {
	var qrCodePrefix sql.NullString
	var cmcId sql.NullInt64
	var parentSymbol sql.NullString

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
		&currConf.ModuleType,
		&currConf.Active,
		&currConf.LastUpdated,
	)
	if err != nil {
		return errs.AddTrace(err)
	}

	if qrCodePrefix.Valid {
		currConf.QrCodePrefix = qrCodePrefix.String
	}
	if cmcId.Valid {
		currConf.CmcId = int(cmcId.Int64)
	}
	if parentSymbol.Valid {
		currConf.ParentSymbol = parentSymbol.String
	}

	return nil
}

func (r *currencyConfigRepository) Update(currencyConfig cc.CurrencyConfig) (err error) {
	err = r.db.QueryRow(`
    UPDATE `+currencyConfigTable+`
    SET 
        symbol = ?,
        name = ?,
        unit = ?,
        token_type = ?,
        is_finance_enabled = ?,
        is_single_address = ?,
        is_using_memo = ?,
        is_qrcode_enabled = ?,
        is_address_notice_enabled = ?,
        qrcode_prefix = ?,
        withdraw_fee = ?,
        healthy_block_diff = ?,
        default_idr_price = ?,
        cmc_id = ?,
        parent_symbol = ?,
        module_type = ?,
        last_updated = now()
    WHERE id = ?`,
		currencyConfig.Symbol,
		currencyConfig.Name,
		currencyConfig.Unit,
		currencyConfig.TokenType,
		currencyConfig.IsFinanceEnabled,
		currencyConfig.IsSingleAddress,
		currencyConfig.IsUsingMemo,
		currencyConfig.IsQrCodeEnabled,
		currencyConfig.IsAddressNoticeEnabled,
		currencyConfig.QrCodePrefix,
		currencyConfig.WithdrawFee,
		currencyConfig.HealthyBlockDiff,
		currencyConfig.DefaultIdrPrice,
		currencyConfig.CmcId,
		currencyConfig.ParentSymbol,
		currencyConfig.ModuleType,
		currencyConfig.Id,
	).Err()
	if err != nil {
		return errs.AddTrace(err)
	}
	return
}

func (r *currencyConfigRepository) ToggleActive(userId int, active bool) error {
	query := "UPDATE " + currencyConfigTable + " SET active = ? WHERE id = ?"

	err := r.db.QueryRow(query, active, userId).Err()
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
