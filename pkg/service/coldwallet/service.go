package coldwallet

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	handlers "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type ColdWalletService interface {
	ActivateColdWallet(ctx context.Context, id int) (err error)
	DeactivateColdWallet(ctx context.Context, id int) (err error)
	ListColdWallet(ctx context.Context, page int, limit int) (resp []cb.ColdBalance, err error)
	GetBalance(ctx context.Context, currencyConfigId int) (coldBalances []cb.ColdBalance)
	CreateColdWallet(ctx context.Context, createReq domain.CreateColdBalance) (err error)
	UpdateColdWallet(ctx context.Context, updateReq domain.ColdBalance) (err error)
	SendToHot(ctx context.Context, sendToHotReq handlers.SendToHotReq) (err error)
	SettlementWallet(ctx context.Context, currencyId int) (result domain.ColdBalance, err error)
	UpdateBalance(ctx context.Context, id int, balance string) (err error)
}

type coldWalletService struct {
	validator util.CustomValidator
	cbRepo    domain.Repository
}

func NewColdWalletService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *coldWalletService {
	return &coldWalletService{
		validator,
		mysqlRepos.ColdBalance,
	}
}
