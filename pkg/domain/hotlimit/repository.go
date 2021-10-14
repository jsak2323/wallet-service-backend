package hotlimit

type Repository interface {
	GetBySymbol(symbol string) (HotLimit, error)
}