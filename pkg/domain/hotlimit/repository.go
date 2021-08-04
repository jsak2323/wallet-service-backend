package hotlimit

type Repository interface {
	GetByCurrencyId(currencyId int) (map[string]HotLimit, error)
}