package currency

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	cr "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	currencyHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/currency"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type CurrencyService interface {
	ActivateCurrency(ctx context.Context, id int) (err error)
	DeactivateCurrencyConfig(ctx context.Context, userId int) (err error)
	CreateCurrencyConfig(ctx context.Context, req domain.CurrencyConfig) (err error)
	UpdateCurrencyConfig(ctx context.Context, currencyConfig domain.UpdateCurrencyConfig) (err error)
	CreateCurrencyRpc(ctx context.Context, cRreq currencyHandler.CurrencyRpcReq) (err error)
	ListCurrency(ctx context.Context) (RES currencyHandler.ListRes, err error)
	DeleteCurrencyRpc(ctx context.Context, currencyId int, rpcId int) (err error)
}

type currencyService struct {
	validator util.CustomValidator
	ccRepo    cc.Repository
	crRepo    cr.Repository
	rcRepo    rc.Repository
}

func NewCurrencyService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *currencyService {
	return &currencyService{
		validator,
		mysqlRepos.CurrencyConfig,
		mysqlRepos.CurrencyRpc,
		mysqlRepos.RpcConfig,
	}
}
