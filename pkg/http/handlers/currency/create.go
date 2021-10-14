package currency

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
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
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- currency.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully created"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.CreateHandler For all symbols, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.CreateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.ccRepo.Create(currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.CreateHandler ccRepo.GetAll Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(currencyConfig domain.CurrencyConfig) error {
	if currencyConfig.Symbol == "" {
		return errors.New("Symbol")
	}
	if currencyConfig.Name == "" {
		return errors.New("Name")
	}
	if currencyConfig.Unit == "" {
		return errors.New("Unit")
	}

	return nil
}
