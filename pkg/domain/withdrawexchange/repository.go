package withdrawexchange

type Repository interface {
	GetPendingWithdraw(symbol string) (string, error)
}