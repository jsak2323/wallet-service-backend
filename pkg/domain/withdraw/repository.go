package withdraw

type Repository interface {
	GetPendingWithdraw(symbol string) (string, error)
}