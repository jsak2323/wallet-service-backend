package mysql

// import (
// 	"database/sql"
	
// 	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
// )

// const hotLimitTable = "hot_limit"
// const hotLimitDefaultLimit = 10

// type hotLimitRepository struct {
// 	db *sql.DB
// }

// func NewMysqlHotLimitRepository(db *sql.DB) domain.Repository {
//     return &hotLimitRepository{
//         db,
//     }
// }

// func(r *hotLimitRepository) GetByCurrencyId(currencyId int) (limits map[string]domain.HotLimit, err error) {
// 	query := "SELECT id, currency_id, type, amount FROM " + hotLimitTable + " WHERE currency_id = ?"
	
// 	return r.queryRows(query, currencyId)
// }

// func (r *hotLimitRepository) queryRows(query string, params... interface{}) (limits map[string]domain.HotLimit, err error) {
// 	limits = make(map[string]domain.HotLimit, 5)
	
// 	rows, err := r.db.Query(query, params...)
// 	if err != nil {
// 		return map[string]domain.HotLimit{}, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var limit domain.HotLimit
		
// 		if err = rows.Scan(
// 			&limit.Id,
// 			&limit.CurrencyId,
// 			&limit.Type,
// 			&limit.Amount,
// 		); err != nil {
// 			return map[string]domain.HotLimit{}, err
// 		}

// 		limits[limit.Type] = limit
// 	}

// 	return limits, nil
// }