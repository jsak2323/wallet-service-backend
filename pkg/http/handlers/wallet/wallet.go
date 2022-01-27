package wallet

import (
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	w "github.com/btcid/wallet-services-backend-go/pkg/domain/withdrawexchange"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	modules "github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type WalletService struct {
	moduleServices       *modules.ModuleServiceMap
	coldWalletService    *hcw.ColdWalletService
	marketService        *h.MarketService
	withdrawExchangeRepo w.Repository
	hotLimitRepo         hl.Repository
	userBalanceRepo      ub.Repository
}

func NewWalletService(
	moduleServices *modules.ModuleServiceMap,
	coldWalletService *hcw.ColdWalletService,
	marketService *h.MarketService,
	withdrawExchangeRepo w.Repository,
	hotLimitRepo hl.Repository,
	userBalanceRepo ub.Repository,
) *WalletService {
	return &WalletService{
		moduleServices, coldWalletService, marketService, withdrawExchangeRepo, hotLimitRepo, userBalanceRepo,
	}
}
