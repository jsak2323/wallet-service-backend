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

func (s *RpcConfigService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcConfig domain.RpcConfig
		RES       StandardRes
		err       error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcconfig.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = validateCreateReq(rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.rcRepo.Create(rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfig)
		return
	}
}

func validateCreateReq(rpcConfig domain.RpcConfig) error {
	if rpcConfig.Name == "" {
		return errs.AddTrace(errors.New("invalid Name"))
	}
	if rpcConfig.Host == "" {
		return errs.AddTrace(errors.New("invalid Host"))
	}
	if rpcConfig.Path == "" {
		return errs.AddTrace(errors.New("invalid Path"))
	}

	return nil
}
