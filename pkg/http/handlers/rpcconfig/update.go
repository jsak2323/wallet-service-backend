package rpcconfig

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
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
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcconfig.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateUpdateReq(rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.UpdateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rcRepo.Update(rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.UpdateHandler rcRepo.Update Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(rpcConfig domain.RpcConfig) error {
	if rpcConfig.Id == 0 {
		return errors.New("ID")
	}
	if rpcConfig.CurrencyId == 0 {
		return errors.New("Currency ID")
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
