package walletrpc

import (
	"context"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	handler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type WalletRpcService interface {
	GetAddressType(ctx context.Context, symbol string, tokenType string, address string) (resp handler.AddressTypeRes, err error)
	InvokeGetBalance(ctx context.Context, symbol, tokenType string) (RES *GetBalanceHandlerResponseMap)
	InvokeGetBlockCount(ctx context.Context, symbol, tokenType string) (RES *GetBlockCountHandlerResponseMap, err error)
	InvokeGetHealthCheck(ctx context.Context, symbol, tokenType string) (RES *GetHealthCheckHandlerResponseMap, err error)
	GetLog(ctx context.Context, SYMBOL string, TOKENTYPE string, rpcConfigType string, date string) (resp *http.Response, err error)
	GetNewAddress(ctx context.Context, symbol string, tokenType string, addressType string) (resp *handlers.GetNewAddressRes, err error)
	InvokeListTransactions(ctx context.Context, symbol, tokenType string, limit int) (RES *ListTransactionsHandlerResponseMap, err error)
	InvokeListWithdraws(ctx context.Context, symbol, tokenType string, limit int) (RES *ListWithdrawsHandlerResponseMap, err error)
	SendToAddress(ctx context.Context, req handlers.SendToAddressRequest) (resp *handlers.SendToAddressRes, err error)
}

type walletRpcService struct {
	validator        util.CustomValidator
	moduleServices   *modules.ModuleServiceMap
	systemConfigRepo sc.Repository
	healthCheckRepo  hc.Repository
}

func NewWalletRpcService(validator util.CustomValidator, moduleServices *modules.ModuleServiceMap, mysqlRepos mysql.MysqlRepositories) *walletRpcService {
	return &walletRpcService{
		validator,
		moduleServices,
		mysqlRepos.SystemConfig,
		mysqlRepos.HealthCheck,
	}
}
