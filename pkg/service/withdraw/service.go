package withdraw

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/withdraw"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type WithdrawService interface {
	List(ctx context.Context, page int, limit int, filters []map[string]interface{}) (resp []withdraw.Withdraw, err error)
}

type withdrawService struct {
	validator    util.CustomValidator
	withdrawRepo withdraw.Repository
}

func NewWithdrawService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *withdrawService {
	return &withdrawService{
		validator,
		mysqlRepos.Withdraw,
	}
}
