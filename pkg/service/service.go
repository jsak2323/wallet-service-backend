package service

import (
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/service/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type Service struct {
	Permission permission.PermissionService
}

func New(
	validator util.CustomValidator,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
) Service {
	svc := Service{
		Permission: permission.NewPermissionService(validator, mysqlRepos, exchangeApiRepos),
	}
	return svc
}
