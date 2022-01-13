package currency

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		currencyConfig domain.CurrencyConfig
		RES            StandardRes
		err            error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			RES.Success = true
			RES.Message = "Currency successfully created"

			config.LoadCurrencyConfigs()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.CreateHandler For all symbols, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = validateCreateReq(currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.ccRepo.Create(currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateCurrencyConfig)
		return
	}
}

func validateCreateReq(currencyConfig domain.CurrencyConfig) error {
	if currencyConfig.Symbol == "" {
		return errs.AddTrace(errors.New("invalid symbol"))
	}
	if currencyConfig.Name == "" {
		return errs.AddTrace(errors.New("invalid name"))
	}
	if currencyConfig.Unit == "" {
		return errs.AddTrace(errors.New("invalid unit"))
	}

	return nil
}
