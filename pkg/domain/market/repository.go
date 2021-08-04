package market

type Repository interface {
	LastPriceBySymbol(symbol string) (string, error)
}