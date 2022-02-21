package market

type Repository interface {
	LastPriceBySymbol(symbol, trade string) (string, error)
}
