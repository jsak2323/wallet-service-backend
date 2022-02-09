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
		currencyConfig domain.UpdateCurrencyConfig
		RES            StandardRes
		err            error
		ctx            = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- currency.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs(ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- currency.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = s.validator.Validate(currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.ccRepo.Update(ctx, currencyConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateCurrencyConfig)
		return
	}
}

func validateUpdateReq(currencyConfig domain.CurrencyConfig) error {
	if currencyConfig.Id == 0 {
		return errs.AddTrace(errors.New("invalid ID"))
	}
	if currencyConfig.Symbol == "" {
		return errs.AddTrace(errors.New("invalid Symbol"))
	}
	if currencyConfig.Name == "" {
		return errs.AddTrace(errors.New("invalid Name"))
	}
	if currencyConfig.Unit == "" {
		return errs.AddTrace(errors.New("invalid Unit"))
	}

	return nil
}
