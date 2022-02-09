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
		rpcConfig domain.UpdateRpcConfig
		RES       StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcconfig.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs(ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = s.validator.Validate(rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.rcRepo.Update(ctx, rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCConfig)
		return
	}
}

func validateUpdateReq(rpcConfig domain.RpcConfig) error {
	if rpcConfig.Id == 0 {
		return errs.AddTrace(errors.New("ID"))
	}
	if rpcConfig.Name == "" {
		return errs.AddTrace(errors.New("Name"))
	}
	if rpcConfig.Host == "" {
		return errs.AddTrace(errors.New("Host"))
	}
	if rpcConfig.Path == "" {
		return errs.AddTrace(errors.New("Path"))
	}

	return nil
}
