package wallet

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	w "github.com/btcid/wallet-services-backend-go/pkg/domain/withdrawexchange"
	modules "github.com/btcid/wallet-services-backend-go/pkg/modules"
	cws "github.com/btcid/wallet-services-backend-go/pkg/service/coldwallet"
	ms "github.com/btcid/wallet-services-backend-go/pkg/service/market"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type WalletService interface {
	FormatWalletBalanceCurrency(walletBalance GetBalanceRes) (result GetBalanceRes)
	GetBalance(ctx context.Context, currConfig config.CurrencyRpcConfig) GetBalanceRes
	InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, currencyId int)
	SetColdBalanceDetails(ctx context.Context, res *GetBalanceRes)
	SetHotBalanceDetails(ctx context.Context, rpcConfigs []rc.RpcConfig, res *GetBalanceRes)
	SetHotLimits(ctx context.Context, res *GetBalanceRes)
	SetPendingWithdraw(ctx context.Context, res *GetBalanceRes)
	SetPercent(ctx context.Context, res *GetBalanceRes)
	SetUserBalanceDetails(ctx context.Context, res *GetBalanceRes)
}

type walletService struct {
	moduleServices       *modules.ModuleServiceMap
	coldWalletService    cws.ColdWalletService
	marketService        ms.MarketService
	coldBalance          cb.Repository
	hotLimitRepo         hl.Repository
	userBalanceRepo      ub.Repository
	marketRepo           market.Repository
	withdrawExchangeRepo w.Repository
}

func NewWalletService(
	moduleServices *modules.ModuleServiceMap,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
	coldWalletService cws.ColdWalletService,
	marketService ms.MarketService,
) *walletService {
	return &walletService{
		moduleServices,
		coldWalletService,
		marketService,
		mysqlRepos.ColdBalance,
		exchangeApiRepos.HotLimit,
		mysqlRepos.UserBalance,
		exchangeApiRepos.Market,
		mysqlRepos.WithdrawExchange,
	}
}
