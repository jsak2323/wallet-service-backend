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

func (s *CurrencyConfigService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		currencyConfig domain.CurrencyConfig
		RES            StandardRes
		err            error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- currency.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateUpdateReq(currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.ccRepo.Update(currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateCurrencyConfig.Title})
		return
	}
}

func validateUpdateReq(currencyConfig domain.CurrencyConfig) error {
	if currencyConfig.Id == 0 {
		return errors.New("invalid ID")
	}
	if currencyConfig.Symbol == "" {
		return errors.New("invalid Symbol")
	}
	if currencyConfig.Name == "" {
		return errors.New("invalid Name")
	}
	if currencyConfig.Unit == "" {
		return errors.New("invalid Unit")
	}

	return nil
}
