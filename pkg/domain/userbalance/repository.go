package userbalance

type Repository interface {
	GetTotalCoinBalance(coin string) (TotalCoinBalance, error)
}