package rpcconfig

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcConfig domain.RpcConfig
		RES       StandardRes
		err       error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcconfig.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateUpdateReq(rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.rcRepo.Update(rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateRPCConfig.Title})
		return
	}
}

func validateUpdateReq(rpcConfig domain.RpcConfig) error {
	if rpcConfig.Id == 0 {
		return errors.New("ID")
	}
	if rpcConfig.Name == "" {
		return errors.New("Name")
	}
	if rpcConfig.Host == "" {
		return errors.New("Host")
	}
	if rpcConfig.Path == "" {
		return errors.New("Path")
	}

	return nil
}
