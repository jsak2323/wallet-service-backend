package userwallet

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type UserWalletService interface {
	InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, symbol string)
}

type userWalletService struct {
	userBalanceRepo ub.Repository
}

func NewUserWalet(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *userWalletService {
	return &userWalletService{
		mysqlRepos.UserBalance,
	}
}
