package currency

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *CurrencyConfigService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- currency.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateUpdateReq(currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.UpdateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.ccRepo.Update(currencyConfig); err != nil {
		logger.ErrorLog(" -- currency.UpdateHandler ccRepo.Update Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(currencyConfig domain.CurrencyConfig) error {
	if currencyConfig.Id == 0 {
		return errors.New("ID")
	}
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
