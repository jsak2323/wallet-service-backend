package currency

import (
	"net/http"
    "encoding/json"

    cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	"github.com/btcid/wallet-services-backend-go/cmd/config"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type CurrencyConfigService struct {
	ccRepo cc.CurrencyConfigRepository
}

func NewCurrencyConfigService(ccRepo cc.CurrencyConfigRepository) *CurrencyConfigService {
    return &CurrencyConfigService{ccRepo: ccRepo}
}

func (s CurrencyConfigService) GetCurrencyConfigHandler(w http.ResponseWriter, req *http.Request) { 
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
	
	logger.InfoLog(" - GetCurrencyConfigHandler For all symbols, Requesting ...", req) 

	if len(config.CURR) > 0 {
		for _, curr := range config.CURR {
			RES.CurrencyConfigs = append(RES.CurrencyConfigs, curr.Config)
		}
	} else {
		ccs, err := s.ccRepo.GetAll()
		if err != nil {
			logger.ErrorLog(" -- GetCurrencyConfigHandler ccRepo.GetAll Error: "+err.Error())
			RES.Error = err.Error()
			return
		}

		for _, currConfig := range ccs {
			RES.CurrencyConfigs = append(RES.CurrencyConfigs, currConfig)
		}
	}
}