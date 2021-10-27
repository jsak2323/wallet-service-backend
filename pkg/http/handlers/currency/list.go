package currency

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) ListHandler(w http.ResponseWriter, req *http.Request) {
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

	logger.InfoLog(" - ListHandler For all symbols, Requesting ...", req)

	if len(config.CURRRPC) > 0 {
		for _, curr := range config.CURRRPC {
			RES.CurrencyConfigs = append(RES.CurrencyConfigs, curr.Config)
		}
	} else {
		if RES.CurrencyConfigs, err = s.ccRepo.GetAll(); err != nil {
			logger.ErrorLog(" -- ListHandler ccRepo.GetAll Error: " + err.Error())
			RES.Error = err.Error()
			return
		}
	}
}
