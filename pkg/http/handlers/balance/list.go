package balance
 
import (
	"net/http"
	"encoding/json"
	"sync"
	
	"github.com/btcid/wallet-services-backend-go/cmd/config"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
)

func (s *BalanceService) ListBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	for SYMBOL, curr := range config.CURR {
		var wg sync.WaitGroup
		var currenyData WalletBalance = WalletBalance{
			CurrencyConfig: curr.Config,
			WalletBalance: hw.WalletBalance{
				TotalColdCoin: "0", TotalColdIdr: "0",
				TotalNodeCoin: "0", TotalNodeIdr: "0",
				TotalUserCoin: "0", TotalUserIdr: "0",
			},
		}
		
		wg.Add(3)
		go func() { defer wg.Done(); s.walletService.SetColdBalanceDetails(SYMBOL, &currenyData.WalletBalance) }()
		go func() { defer wg.Done(); s.walletService.SetHotBalanceDetails(SYMBOL, curr.RpcConfigs, &currenyData.WalletBalance) }()
		go func() { defer wg.Done(); s.walletService.SetUserBalanceDetails(SYMBOL, &currenyData.WalletBalance) }()
		wg.Wait()

		s.walletService.FormatWalletBalanceCurrency(SYMBOL, &currenyData.WalletBalance)
		
		RES.Balances = append(RES.Balances, currenyData)
	}
}