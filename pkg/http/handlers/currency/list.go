package currency

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - currency.ListHandler For all symbols, Requesting ...", req)

	if len(config.CURRRPC) > 0 {
		for _, curr := range config.CURRRPC {
			RES.CurrencyConfigs = append(RES.CurrencyConfigs, curr.Config)
		}
	} else {
		if RES.CurrencyConfigs, err = s.ccRepo.GetAll(); err != nil {
			RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllCurrencyConfig)
			return
		}
	}

	for i, currency := range RES.CurrencyConfigs {
		RES.CurrencyConfigs[i].RpcConfigs, err = s.rcRepo.GetByCurrencyId(currency.Id)
		if err != nil {
			RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRPCConfigByCurrencyID)
			return
		}
	}
}
