package wallet

import (
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	modules "github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type WalletService struct {
	moduleServices  	*modules.ModuleServiceMap
    coldWalletService 	hcw.ColdWalletService
	marketService 		h.MarketService
    userBalanceRepo 	ub.Repository
}

func NewWalletService(
	moduleServices  	*modules.ModuleServiceMap,
    coldWalletService 	hcw.ColdWalletService,
	marketService 		h.MarketService,
    userBalanceRepo 	ub.Repository,
) *WalletService {
	return &WalletService{
		moduleServices, coldWalletService, marketService, userBalanceRepo,
	}
}
