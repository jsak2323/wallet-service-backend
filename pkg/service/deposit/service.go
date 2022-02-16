package deposit

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/deposit"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type DepositService interface {
	List(ctx context.Context, page int, limit int, filters []map[string]interface{}) (resps []deposit.Deposit, err error)
}

type depositService struct {
	validator   util.CustomValidator
	depositRepo deposit.Repository
}

func NewDepositService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *depositService {
	return &depositService{
		validator,
		mysqlRepos.Deposit,
	}
}
