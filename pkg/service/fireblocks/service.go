package fireblocks

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type FireblocksService interface {
	ValidateHotDestAddress(ctx context.Context, req fireblocks.FireblocksSignReq) (resp fireblocks.FireblocksSignRes, err error)
}

type fireblocksService struct {
	validator   util.CustomValidator
	coldbalance coldbalance.Repository
}

func NewFireblocksService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *fireblocksService {
	return &fireblocksService{
		validator,
		mysqlRepos.ColdBalance,
	}
}
