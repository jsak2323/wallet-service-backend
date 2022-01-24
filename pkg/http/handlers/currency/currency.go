package currency

import (
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	cr "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyrpc"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

const errInternalServer = "Internal server error"

type CurrencyConfigService struct {
	ccRepo    cc.Repository
	crRepo    cr.Repository
	rcRepo    rc.Repository
	validator util.CustomValidator
}

func NewCurrencyConfigService(
	ccRepo cc.Repository,
	crRepo cr.Repository,
	rcRepo rc.Repository,
	validator util.CustomValidator,
) *CurrencyConfigService {
	return &CurrencyConfigService{
		ccRepo:    ccRepo,
		crRepo:    crRepo,
		rcRepo:    rcRepo,
		validator: validator,
	}
}
