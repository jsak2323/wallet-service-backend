package rpcconfig

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) CreateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- rpcconfig.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.CreateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rcRepo.Create(rpcConfig); err != nil {
		logger.ErrorLog(" -- rpcconfig.CreateHandler rcRepo.Create Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(rpcConfig domain.RpcConfig) error {
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
